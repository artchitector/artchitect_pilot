package entropy

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"image/png"
	"math"
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

func (gk *Gatekeeper) NotifyEntropyPhase(
	ctx context.Context,
	img image.Image,
	phase string,
	entropyAnswer uint64,
) error {
	buf := new(bytes.Buffer)
	var format string
	switch phase {
	case PhaseSource:
		fallthrough
	case PhaseNoise:
		format = model.ImageTypeJPEG
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityM}); err != nil {
			return errors.Wrapf(err, "[gatekeeper] failed to make jpeg")
		}
	case PhaseNoiseShrink:
		fallthrough
	case PhaseBytes:
		format = model.ImageTypePNG
		if err := png.Encode(buf, img); err != nil {
			return errors.Wrapf(err, "[gatekeeper] failed to make png")
		}
	}

	imgBase64Encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	var entropyAnswerByte string
	var entropyAnswerFloat float64

	if entropyAnswer > 0 {
		entropyAnswerByte = fmt.Sprintf("%064b", entropyAnswer)
		entropyAnswerFloat = float64(entropyAnswer) / float64(math.MaxUint64)
	}

	// Теперь отправить уведомление в канал
	err := gk.notifier.NotifyEntropy(ctx, model.EntropyState{
		Timestamp:         time.Now(),
		Phase:             phase,
		Image:             imgBase64Encoded,
		ImageType:         format,
		EntropyAnswerByte: entropyAnswerByte,
		EntropyAnswer:     entropyAnswerFloat,
	})
	if err != nil {
		return errors.Wrapf(err, "[gatekeeper] failed to notify phase")
	}
	return nil
}
