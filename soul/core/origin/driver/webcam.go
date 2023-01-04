package driver

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/bits"
	"net/http"
)

/*
WebcamDriver gets image from webcam (from http-url, where you can get jpg-image from webcam).
Then image normalized into float64 value from 0.0 to 1.0. And this is the answer for all questions.
*/
type WebcamDriver struct {
	originUrl          string
	decisionRepository decisionRepository
}

func NewWebcamDriver(originUrl string, decisionRepository decisionRepository) *WebcamDriver {
	return &WebcamDriver{originUrl, decisionRepository}
}

type decisionRepository interface {
	SaveDecision(ctx context.Context, decision model.Decision) (model.Decision, error)
}

func (w *WebcamDriver) GetValue(ctx context.Context) (float64, error) {
	response, err := http.Get(w.originUrl)
	if err != nil {
		return 0.0, errors.Wrapf(err, "failed to get %s", w.originUrl)
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	//jpeg.Decode()
	if err != nil {
		return 0.0, errors.Wrap(err, "failed to decode image from response.Body")
	}

	return w.imageToNumber(ctx, img)
}

func (w *WebcamDriver) imageToNumber(ctx context.Context, img image.Image) (float64, error) {
	img = resize.Resize(4, 2, img, resize.Lanczos3)

	size := img.Bounds().Size()
	bts := make([]uint8, 0, size.X*size.Y)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			red := float64(originalColor.R)
			green := float64(originalColor.G)
			blue := float64(originalColor.B)

			grey := uint8(math.Round((red + green + blue) / 3))
			bts = append(bts, grey)
		}
	}

	// another strategy of random generation. I don't know, what is best.
	//result := uint64(0)
	//for idx, bt := range bts {
	//	// each 8bit number reverts. this hack gives more "random" results.
	//	mask := uint64(bits.Reverse8(bt))
	//	// integers shifted into it's place in 64-bit map
	//	maskShifted := mask << ((len(bts) - 1 - idx) * 8)
	//	// result is a combination of 8 bitmasks (size of 8bit), reversed and shifted.
	//	result = result | maskShifted
	//}

	result := uint64(0)
	for idx, bt := range bts {
		mask := uint64(bt)
		// integers shifted into it's place in 64-bit map
		maskShifted := mask << ((len(bts) - 1 - idx) * 8)
		// each shifted uint64 will be reversed. First pixel will be in the end of chain.
		maskReversed := bits.Reverse64(maskShifted)
		// result is a combination of 8 bitmasks (size of 8bit), shifted and reversed.
		result = result | maskReversed
	}

	flResult := float64(result) / float64(math.MaxUint64)
	log.Info().Msgf("[webcam] generated number: %d. Meaning: %f", result, flResult)
	go func() {
		if err := w.saveDecision(ctx, img, flResult); err != nil {
			log.Error().Err(err).Msg("[webcam] failed to save decision")
		}
	}()
	return flResult, nil
}

func (w *WebcamDriver) saveDecision(ctx context.Context, img image.Image, result float64) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return errors.Wrap(err, "failed to encode image to jpeg")
	}
	if decision, err := w.decisionRepository.SaveDecision(ctx, model.Decision{
		Output:              result,
		Artifact:            buf.Bytes(),
		ArtifactContentType: model.ArtifactContentTypeJpeg,
	}); err != nil {
		return errors.Wrapf(err, "failed to save decision with result=%f", result)
	} else {
		log.Info().Msgf("[webcam] save decision id=%d with result=%f", decision.ID, decision.Output)
	}

	return nil
}
