package driver

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
)

type changeable interface {
	Set(x, y int, c color.Color)
}

/*
WebcamDriver gets image from webcam (from http-url, where you can get jpg-image from webcam).
Then image normalized into float64 value from 0.0 to 1.0. And this is the answer for all questions.
*/
type WebcamDriver struct {
	originUrl string
	lastFrame image.Image
}

func NewWebcamDriver(originUrl string) *WebcamDriver {
	return &WebcamDriver{originUrl, nil}
}

func (w *WebcamDriver) GetValue(ctx context.Context) (float64, error) {
	response, err := http.Get(w.originUrl)
	if err != nil {
		return 0.0, errors.Wrapf(err, "failed to get %s", w.originUrl)
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return 0.0, errors.Wrap(err, "failed to decode image from response.Body")
	}

	result, err := w.imageToNumber(ctx, img)
	log.Info().Msgf("[webcam] got number %f", result)
	return result, err
}

func (w *WebcamDriver) imageToNumber(ctx context.Context, originalImg image.Image) (float64, error) {
	var result float64
	var err error
	result, err = w.imageToNumberHash(ctx, originalImg)
	if err != nil {
		return 0.0, errors.Wrapf(err, "[webcam] failed to imageToNumberHash")
	}
	//
	//// testing random generation with entropy usage
	//if entropyResult, err := w.testNewEntropyReading(ctx, originalImg); err != nil {
	//	log.Error().Err(err).Send()
	//} else {
	//	log.Info().Msgf("[ENTROPY] result=%.12f", entropyResult)
	//}

	return result, nil
}

func (w *WebcamDriver) imageToNumberHash(ctx context.Context, originalImg image.Image) (float64, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, originalImg); err != nil {
		return 0.0, errors.Wrap(err, "failed to encode png to bytes")
	}

	hash := md5.Sum(buf.Bytes())

	var result uint
	for _, b := range hash[:8] {
		result = (result << 8) | uint(b)
	}

	flResult := float64(result) / float64(math.MaxUint)
	log.Debug().Msgf("[webcam][imageToNumberHash] generated number: %d. Meaning: %.12f", result, flResult)
	return flResult, nil
}

//
//func (w *WebcamDriver) testNewEntropyReading(ctx context.Context, originalImg image.Image) (float64, error) {
//	/*
//		Алгоритм преобразования:
//		1) берётся 2 соседних кадра
//		2) считается только разница между этими кадрами, она становится новой картинкой (пиксели вычитаются)
//		3) разница усиливается
//		4) картинка сжимается до малого количества пикселей
//		5) из этой картинки снимается float-число (как, пока непонятно)
//	*/
//	if w.lastFrame == nil {
//		w.lastFrame = originalImg
//		return 0.0, errors.New("[ENTROPY] no lastFrame to calculate difference")
//	}
//	clr := originalImg.At(0, 0)
//	log.Info().Msgf("[ENTROPY] color: %+v", clr)
//
//	diffImg := image.NewRGBA(originalImg.Bounds())
//
//	bounds := originalImg.Bounds()
//	for x := 0; x <= bounds.Dx(); x++ {
//		for y := 0; y <= bounds.Dy(); y++ {
//			aClr := w.lastFrame.At(x, y)
//			bClr := originalImg.At(x, y)
//
//			if _, ok := aClr.(color.RGBA); !ok {
//				return 0.0, errors.New("[ENTROPY] A not RGBA color")
//			}
//			if _, ok := bClr.(color.RGBA); !ok {
//				return 0.0, errors.New("[ENTROPY] B not RGBA color")
//			}
//
//			var newR, newG, newB int16
//			newR = newR + int16(bClr.(color.RGBA).R) - int16(aClr.(color.RGBA).R)
//			if newR < 0 {
//				newR *= -1
//			}
//			newG = newG + int16(bClr.(color.RGBA).G) - int16(aClr.(color.RGBA).G)
//			if newG < 0 {
//				newG *= -1
//			}
//			newB = newB + int16(bClr.(color.RGBA).B) - int16(aClr.(color.RGBA).B)
//			if newB < 0 {
//				newB *= -1
//			}
//
//			newColor := color.RGBA{
//				R: uint8(newR * 10),
//				G: uint8(newG * 10),
//				B: uint8(newB * 10),
//				A: 255,
//			}
//			diffImg.SetRGBA(x, y, newColor)
//		}
//	}
//
//	w.saveImg(ctx, diffImg, "diff1")
//
//	return 0.0, nil
//}
//
//func (w *WebcamDriver) saveImg(ctx context.Context, img image.Image, name string) {
//	buf := new(bytes.Buffer)
//	if err := png.Encode(buf, img); err != nil {
//		log.Error().Err(err).Send()
//		return
//	}
//	if err := os.WriteFile(name+".png", buf.Bytes(), 0777); err != nil {
//		log.Error().Err(err).Send()
//		return
//	}
//	log.Info().Msgf("[FILE] saved file %s", name)
//}
