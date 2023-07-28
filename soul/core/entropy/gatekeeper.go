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
	"image/png"
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

func (gk *Gatekeeper) NotifyEntropyState(
	ctx context.Context,
	state model.EntropyState,
) error {
	// convert images to base64 form
	for key, img := range state.Images {
		buf := new(bytes.Buffer)
		if key == ImageSource || key == ImageNoise {
			if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityM}); err != nil {
				return errors.Wrapf(err, "[gatekeeper] failed to encode jpeg key=%s", key)
			}
		} else {
			if err := png.Encode(buf, img); err != nil {
				return errors.Wrapf(err, "[gatekeeper] failed to encode png key=%s", key)
			}
		}
		imgBase64Encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
		state.ImagesEncoded[key] = imgBase64Encoded
	}
	state.Images = make(map[string]image.Image) // clear images
	err := gk.notifier.NotifyEntropy(ctx, state)
	if err != nil {
		return errors.Wrapf(err, "[gatekeeper] failed to notify phase")
	}
	return nil
}
