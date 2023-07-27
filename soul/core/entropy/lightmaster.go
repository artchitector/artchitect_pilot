package entropy

import (
	"context"
	"github.com/artchitector/artchitect/soul/resources"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/color"
	"math"
)

const (
	PhaseSource      = "source"
	PhaseNoise       = "noise"
	PhaseNoiseShrink = "shrink"
	PhaseBytes       = "bytes"
	LastFramesToUse  = 2 // Шум считается между двумя или более кадров
	SquareSize       = 64 * 7
	ResultSize       = 8
)

/*
Lightmaster отслеживает состояние энтропии в текущем кадре
Еще он передаёт детализацию обработки энтропии на gate-сервер через redis. Это нужно, чтобы на клиенте был виден
постоянный процесс обработки энтропии в виде jpeg-стримов, и видно было как картинка превращается в решение.

Не каждое состояние используется в принятии решений, многие пропускаются.
*/
type Lightmaster struct {
	webcam      *resources.Webcam
	gatekeeper  *Gatekeeper
	lastNFrames []image.Image
}

func NewLightmaster(webcam *resources.Webcam, gatekeeper *Gatekeeper) *Lightmaster {
	return &Lightmaster{
		webcam,
		gatekeeper,
		make([]image.Image, 0, LastFramesToUse),
	}
}

/*
Запускается процесс считывания кадров с веб-камеры и превращение их в float64-число
*/
func (l *Lightmaster) StartEntropyReading(ctx context.Context) error {
	ch := l.webcam.GetStream(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case i := <-ch:
			if err := l.handleSingleFrame(ctx, i); err != nil {
				log.Error().Err(err).Msgf("[lightmaster] errored handle single frame")
			}
		}
	}
}

func (l *Lightmaster) handleSingleFrame(ctx context.Context, newFrame image.Image) error {
	if err := l.gatekeeper.NotifyEntropyPhase(ctx, newFrame, PhaseSource, 0); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseSource)
		// not stop
	}

	square, err := l.extractSquare(newFrame)
	if err != nil {
		return errors.Wrap(err, "[lightmaster] failed extract square")
	}

	if len(l.lastNFrames) < LastFramesToUse {
		l.lastNFrames = append(l.lastNFrames, square)
	} else {
		l.lastNFrames = append(l.lastNFrames[1:], square)
	}

	if len(l.lastNFrames) == LastFramesToUse {
		if err := l.pipelineEntropy(ctx); err != nil {
			log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoise)
		}
	}

	return nil
}

func (l *Lightmaster) extractSquare(frame image.Image) (image.Image, error) {
	oldBounds := frame.Bounds()
	if oldBounds.Dx() < SquareSize || oldBounds.Dy() < SquareSize {
		return nil, errors.Errorf("[lightmaster] too small image. size is %d and %d", oldBounds.Dx(), oldBounds.Dy())
	}
	squareRect := image.Rect(0, 0, SquareSize, SquareSize)
	squareImg := image.NewRGBA(squareRect)

	leftOffset := (oldBounds.Dx() - squareRect.Dx()) / 2
	topOffset := (oldBounds.Dy() - squareRect.Dy()) / 2

	for x := 0; x < SquareSize; x++ {
		for y := 0; y < SquareSize; y++ {
			squareImg.Set(x, y, frame.At(x+leftOffset, y+topOffset))
		}
	}

	return squareImg, nil
}

func (l *Lightmaster) pipelineEntropy(ctx context.Context) error {
	noiseImage, err := l.sourceToNoise()

	if err != nil {
		return errors.Wrapf(err, "[lightmaster] failed to transform source to noise")
	} else if err := l.gatekeeper.NotifyEntropyPhase(ctx, noiseImage, PhaseNoise, 0); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoise)
	}

	shrinkedNoise, err := l.shrinkNoise(noiseImage)
	if err != nil {
		return errors.Wrapf(err, "[lightmaster] failed to shrink noise")
	} else if err := l.gatekeeper.NotifyEntropyPhase(ctx, shrinkedNoise, PhaseNoiseShrink, 0); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoiseShrink)
	}

	bytesImage, entropyAnswer, err := l.noiseToBytes(shrinkedNoise)
	if err != nil {
		return errors.Wrapf(err, "[lightmaster] failed noise2bytes")
	} else if err := l.gatekeeper.NotifyEntropyPhase(ctx, bytesImage, PhaseBytes, entropyAnswer); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseBytes)
	}

	return nil
}

func (l *Lightmaster) sourceToNoise() (image.Image, error) {
	lastFrame := l.lastNFrames[0] // take most distant frame as A, and current frame as B. B-A=noise
	newFrame := l.lastNFrames[LastFramesToUse-1]
	bounds := newFrame.Bounds()
	noiseImage := image.NewRGBA(bounds)

	for x := 0; x <= bounds.Dx(); x++ {
		for y := 0; y <= bounds.Dy(); y++ {
			oldColor := lastFrame.At(x, y)
			newColor := newFrame.At(x, y)
			if _, ok := oldColor.(color.RGBA); !ok {
				return nil, errors.New("[lightmaster] old is not RGBA color")
			}
			if _, ok := newColor.(color.RGBA); !ok {
				return nil, errors.New("[lightmaster] new is not RGBA color")
			}
			var newR, newG, newB int16
			newR = newR + int16(newColor.(color.RGBA).R) - int16(oldColor.(color.RGBA).R)
			if newR < 0 {
				newR *= -1
			}
			newG = newG + int16(newColor.(color.RGBA).G) - int16(oldColor.(color.RGBA).G)
			if newG < 0 {
				newG *= -1
			}
			newB = newB + int16(newColor.(color.RGBA).B) - int16(oldColor.(color.RGBA).B)
			if newB < 0 {
				newB *= -1
			}

			noiseColor := color.RGBA{
				R: uint8(newR * 10),
				G: uint8(newG * 10),
				B: uint8(newB * 10),
				A: 255,
			}
			noiseImage.SetRGBA(x, y, noiseColor)
		}
	}
	return noiseImage, nil
}

func (l *Lightmaster) shrinkNoise(noiseImage image.Image) (image.Image, error) {
	noiseBounds := noiseImage.Bounds()
	resultBounds := image.Rect(0, 0, ResultSize, ResultSize)
	resultImg := image.NewRGBA(resultBounds)

	// how many source pixels is in result pixed (default 56px into 1px, 448px-side into 8px-side)
	proportion := noiseBounds.Dx() / resultBounds.Dx()
	for x := 0; x < resultBounds.Dx(); x++ {
		for y := 0; y < resultBounds.Dy(); y++ {
			// ищем средний цвет для этого пикселя, ужимая пиксели
			var sumR, sumG, sumB int64
			for nx := x * proportion; nx < x*proportion+proportion; nx++ {
				for ny := y * proportion; ny < y*proportion+proportion; ny++ {
					clr := noiseImage.At(nx, ny)
					sumR += int64(clr.(color.RGBA).R)
					sumG += int64(clr.(color.RGBA).G)
					sumB += int64(clr.(color.RGBA).B)
				}
			}
			sumR /= int64(proportion)
			sumG /= int64(proportion)
			sumB /= int64(proportion)
			resultImg.SetRGBA(x, y, color.RGBA{R: uint8(sumR), G: uint8(sumG), B: uint8(sumB), A: 255})
		}
	}

	return resultImg, nil
}

func (l *Lightmaster) noiseToBytes(noiseImage image.Image) (image.Image, uint64, error) {
	bounds := noiseImage.Bounds()
	bytesImage := image.NewRGBA(bounds)
	var entropyAnswer uint64

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			// считаем сумму заполнения трёх цветов, и делим на три. если это больше 255/2 - цвет бита белый, иначе чёрный
			var sum int16
			clr := noiseImage.At(x, y)
			sum += int16(clr.(color.RGBA).R)
			sum += int16(clr.(color.RGBA).G)
			sum += int16(clr.(color.RGBA).B)
			sum /= 3
			if sum > math.MaxInt8/2 { // set white
				bytesImage.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})

				byteIndex := x*ResultSize + y
				entropyAnswer = entropyAnswer | 1<<(63-byteIndex)
			} else {
				bytesImage.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	return bytesImage, entropyAnswer, nil
}
