package infrastructure

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"sync"
)

type Giver struct {
	input  chan model.Pray
	output chan model.Gift
}

type Prayer struct {
	pray        model.Pray
	takeChannel chan model.Gift
}

/*
Cloud works like bus.
Cloud receives prays, then transfer this prays to serving services, then serving services give gifts
And then prayer receives gift via callback, registered in the cloud.
Working via channels.
First simple version of protocol.

TODO: Need make named prays. Many prayers can make its own named prays and receive only it's answer. Something like RPC.
*/
type Cloud struct {
	mutex   sync.Mutex
	logger  zerolog.Logger
	givers  map[string]Giver
	prayers map[string][]Prayer
}

func NewCloud(logger zerolog.Logger) *Cloud {
	return &Cloud{
		logger:  logger,
		givers:  make(map[string]Giver),
		prayers: make(map[string][]Prayer),
	}
}

func (c *Cloud) Pray(ctx context.Context, pray model.Pray) (chan model.Gift, error) {
	if pray.Name == "" {
		return make(chan model.Gift), errors.New("empty pray name")
	}
	prayer := Prayer{
		pray:        pray,
		takeChannel: make(chan model.Gift),
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, found := c.prayers[pray.Name]
	if !found {
		c.prayers[pray.Name] = []Prayer{prayer}
	} else {
		c.prayers[pray.Name] = append(c.prayers[pray.Name], prayer)
	}

	if err := c.emitPray(ctx, pray); err != nil {
		return prayer.takeChannel, errors.Wrap(err, "failed to emit pray")
	}

	return prayer.takeChannel, nil
}

// Serve - servant-service subscribe to channel with incoming prays, handles it and then give gift via gift-channel
func (c *Cloud) Serve(ctx context.Context, prayName string) (chan model.Pray, chan model.Gift, error) {
	giver := Giver{
		input:  make(chan model.Pray),
		output: make(chan model.Gift),
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, found := c.givers[prayName]
	if found {
		return giver.input, giver.output, errors.Errorf("giver %s already registered in the cloud", prayName)
	}

	c.givers[prayName] = giver

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info().Msgf("shutting down listening of giver takeChannel: %s. ctx.Done", prayName)
				return
			case gift := <-giver.output:
				c.logger.Debug().Msgf("cloud got gift from giver %s", prayName)
				if err := c.deliveryGift(ctx, prayName, gift); err != nil {
					c.logger.Error().Err(err).Msgf("gift delivery failed for giver %s", prayName)
				}
			}
		}
	}()

	return giver.input, giver.output, nil
}

func (c *Cloud) emitPray(ctx context.Context, pray model.Pray) error {
	giver, found := c.givers[pray.Name]
	if !found {
		return errors.Errorf("no giver for prays %s", pray.Name)
	}
	go func() {
		giver.input <- pray
	}()
	return nil
}

func (c *Cloud) deliveryGift(ctx context.Context, prayName string, gift model.Gift) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	waitingPrayers, found := c.prayers[prayName]
	c.mutex.Unlock()

	if !found {
		c.logger.Warn().Msgf("no prayers to take gift %s", prayName)
		return nil
	}
	for _, prayer := range waitingPrayers {
		go func(prayer Prayer) { prayer.takeChannel <- gift }(prayer)
	}
	c.mutex.Lock() // unlock in defer section
	c.prayers[prayName] = make([]Prayer, 0)
	return nil
}
