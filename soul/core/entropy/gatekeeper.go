package entropy

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"time"
)

type notifier interface {
	NotifyEntropy(ctx context.Context, state model.EntropyState) error
}

/*
Gatekeeper отправяет сообщения на гейт с изображениями разных фаз работы с энтропией
*/
type Gatekeeper struct {
	redises  map[string]*redis.Client
	notifier notifier
}

func NewGatekeeper(redises map[string]*redis.Client, notifier notifier) *Gatekeeper {
	return &Gatekeeper{redises, notifier}
}

func (gk *Gatekeeper) PushPhase(ctx context.Context, img image.Image, phase string) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityS}); err != nil {
		return errors.Wrapf(err, "[gatekeeper] failed to make jpeg")
	}

	imgBase64Encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Теперь отправить уведомление в канал
	err := gk.notifier.NotifyEntropy(ctx, model.EntropyState{
		Timestamp: time.Now(),
		Phase:     phase,
		Image:     imgBase64Encoded,
	})
	if err != nil {
		return errors.Wrapf(err, "[gatekeeper] failed to notify phase")
	}
	return nil
}
