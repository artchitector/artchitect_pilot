package entropy

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/soul/resources"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/color"
	"math"
	"math/bits"
	"sync"
	"time"
)

const (
	ImageSource  = "source"
	ImageNoise   = "noise"
	ImageEntropy = "entropy"
	ImageChoice  = "choice"

	LastFramesToUse = 2 // Шум считается между двумя или более кадров
	SquareSize      = 64 * 7
	ResultSize      = 8
)

/*
Lightmaster отслеживает состояние энтропии в текущем кадре
Еще он передаёт детализацию обработки энтропии на gate-сервер через redis. Это нужно, чтобы на клиенте был виден
постоянный процесс обработки энтропии в виде jpeg-стримов, и видно было как картинка превращается в решение.

Не каждое состояние используется в принятии решений, многие пропускаются.
*/
type Lightmaster struct {
	webcam        *resources.Webcam
	gatekeeper    *Gatekeeper
	lastNFrames   []image.Image
	tags          []string
	selectedWords map[string]int
	counter       int

	entropyMutex sync.Mutex

	lastEntropyValueUsed bool
	lastEntropyValue     float64
	lastChoiceValueUsed  bool
	lastChoiceValue      float64
}

func NewLightmaster(webcam *resources.Webcam, gatekeeper *Gatekeeper) *Lightmaster {
	return &Lightmaster{
		webcam,
		gatekeeper,
		make([]image.Image, 0, LastFramesToUse),
		nil,
		make(map[string]int),
		0,

		sync.Mutex{},

		false,
		-1.0,
		false,
		-1.0,
	}
}

func (l *Lightmaster) GetEntropy(ctx context.Context) float64 {
	if l.lastEntropyValue >= 0.0 && !l.lastEntropyValueUsed {
		l.entropyMutex.Lock()
		defer l.entropyMutex.Unlock()

		l.lastEntropyValueUsed = true
		return l.lastEntropyValue
	}

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[lightmaster] stop reading entropy, while ctx.Done")
			return 0.0
		case <-time.After(time.Second * 5):
			log.Error().Msgf("[lightmaster] too slow entropy get")
			return 0.0
		case <-time.Tick(time.Millisecond * 50):
			if l.lastEntropyValue >= 0.0 && !l.lastEntropyValueUsed {
				l.entropyMutex.Lock()
				defer l.entropyMutex.Unlock()

				l.lastEntropyValueUsed = true
				return l.lastEntropyValue
			}
		}
	}
}

func (l *Lightmaster) GetChoice(ctx context.Context) float64 {
	if l.lastChoiceValue >= 0.0 && !l.lastChoiceValueUsed {
		l.entropyMutex.Lock()
		defer l.entropyMutex.Unlock()

		l.lastChoiceValueUsed = true
		return l.lastChoiceValue
	}

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[lightmaster] stop reading entropy, while ctx.Done")
			return 0.0
		case <-time.After(time.Second * 5):
			log.Error().Msgf("[lightmaster] too slow entropy get")
			return 0.0
		case <-time.Tick(time.Millisecond * 50):
			if l.lastChoiceValue >= 0.0 && !l.lastChoiceValueUsed {
				l.entropyMutex.Lock()
				defer l.entropyMutex.Unlock()

				l.lastChoiceValueUsed = true
				return l.lastChoiceValue
			}
		}
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
	// step 1. source frame here
	state := model.EntropyState{
		IsShort:       false,
		Images:        make(map[string]image.Image),
		ImagesEncoded: make(map[string]string),
		Entropy:       model.EntropyValue{},
		Choice:        model.EntropyValue{},
	}

	borderedFrame := l.addBordersOnFrame(newFrame)
	state.Images["source"] = borderedFrame

	// step 2. extract square image from source
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
		if err := l.pipelineEntropy(ctx, &state); err != nil {
			return errors.Wrap(err, "[lightmaster] failed to pipeline entropy")
		} else {
			l.entropyMutex.Lock()
			defer l.entropyMutex.Unlock()

			l.lastEntropyValue = state.Entropy.Float64
			l.lastEntropyValueUsed = false
			l.lastChoiceValue = state.Choice.Float64
			l.lastChoiceValueUsed = false
		}
	}

	if err := l.gatekeeper.NotifyEntropyState(ctx, state); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed NotifyEntropyState")
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

func (l *Lightmaster) pipelineEntropy(ctx context.Context, state *model.EntropyState) error {
	noiseImage, err := l.sourceToNoise()

	if err != nil {
		return errors.Wrapf(err, "[lightmaster] failed sourceToNoise")
	} else {
		state.Images[ImageNoise] = noiseImage
	}

	entropyImage, entropyVal, err := l.noiseToEntropy(noiseImage)
	if err != nil {
		return errors.Wrapf(err, "[lightmaster] failed noiseToEntropy")
	} else {
		state.Images[ImageEntropy] = entropyImage
	}
	state.Entropy = l.makeEntropyStruct(entropyVal)

	choiceImage, choiceVal, err := l.invertEntropy(entropyImage)
	if err != nil {
		return errors.Wrapf(err, "[lightmaster] failed invertEntropy")
	} else {
		state.Images[ImageChoice] = choiceImage
	}
	state.Choice = l.makeEntropyStruct(choiceVal)

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

			var amplifierRatio int16 = 30 // чем больше, чем меньше цвета будет на картине
			// ВАЖНО! Усиление и изменение цветов на шумовой картине не влияет на результат. Дальше шум проходит нормализацию,
			// Абсолютные значение не так важны, всё строится на относительной светимости пикселей.
			// Цвет можно выбирать по своему вкусу и дизайну
			noiseColor := color.RGBA{
				R: uint8(255) - uint8(newR*amplifierRatio),
				G: uint8(255) - uint8(newG*amplifierRatio),
				B: uint8(255) - uint8(newB*amplifierRatio),
				A: 255,
			}
			noiseImage.SetRGBA(x, y, noiseColor)
		}
	}
	return noiseImage, nil
}

func (l *Lightmaster) noiseToEntropy(noiseImage image.Image) (image.Image, uint64, error) {
	noiseBounds := noiseImage.Bounds()
	resultBounds := image.Rect(0, 0, ResultSize, ResultSize)
	resultImg := image.NewRGBA(resultBounds)

	proportion := noiseBounds.Dx() / resultBounds.Dx()
	var minPower int64 = math.MaxInt64
	var maxPower int64 = math.MinInt64
	powers := make([][]int64, 0, 8)

	// collect powers of 64-pixels
	for x := 0; x < resultBounds.Dx(); x++ {
		powers = append(powers, make([]int64, 8))
		for y := 0; y < resultBounds.Dy(); y++ {
			var powerOfPixel int64
			// Проходим по всем пикселям в квадрате 56х56 и собираем их силу в сумму
			for nx := x * proportion; nx < x*proportion+proportion; nx++ {
				for ny := y * proportion; ny < y*proportion+proportion; ny++ {
					clr := noiseImage.At(nx, ny)
					powerOfPixel += int64(clr.(color.RGBA).R)
					powerOfPixel += int64(clr.(color.RGBA).G)
					powerOfPixel += int64(clr.(color.RGBA).B)
				}
			}
			if powerOfPixel < minPower {
				minPower = powerOfPixel
			}
			if powerOfPixel > maxPower {
				maxPower = powerOfPixel
			}
			powers[x][y] = powerOfPixel

		}
	}

	var entropyAnswer uint64
	scale := maxPower - minPower
	for x := 0; x < resultBounds.Dx(); x++ {
		for y := 0; y < resultBounds.Dy(); y++ {
			powerOfPixel := powers[x][y] - minPower // clear extra power
			// TODO Если надо инвертировать entropy - Это здесь!
			redPower := math.Round(float64(powerOfPixel) / float64(scale) * 255.0)
			resultImg.SetRGBA(x, y, color.RGBA{R: uint8(redPower), G: uint8(0), B: uint8(0), A: 255})

			if redPower >= 128 {
				byteIndex := x*ResultSize + y
				entropyAnswer = entropyAnswer | 1<<(63-byteIndex)
			}
		}
	}

	return resultImg, entropyAnswer, nil
}

func (l *Lightmaster) invertEntropy(noiseImage image.Image) (image.Image, uint64, error) {
	bounds := noiseImage.Bounds()
	bytesImage := image.NewRGBA(bounds)
	var choiceAnswer uint64

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			// считаем сумму заполнения трёх цветов, и делим на три. если это больше 255/2 - цвет бита белый, иначе чёрный
			clr := noiseImage.At(x, y)
			power := clr.(color.RGBA).R
			power = bits.Reverse8(power)

			if power >= 128 { // set white
				bytesImage.Set(x, y, color.RGBA{R: power, G: 0, B: 0, A: 255})

				byteIndex := x*ResultSize + y
				choiceAnswer = choiceAnswer | 1<<(63-byteIndex)
			} else {
				bytesImage.Set(x, y, color.RGBA{R: power, G: 0, B: 0, A: 255})
			}
		}
	}

	return bytesImage, choiceAnswer, nil
}

func (l *Lightmaster) addBordersOnFrame(frame image.Image) image.Image {
	oldBounds := frame.Bounds()
	if oldBounds.Dx() < SquareSize || oldBounds.Dy() < SquareSize {
		log.Error().Err(errors.Errorf("[lightmaster] too small image. size is %d and %d", oldBounds.Dx(), oldBounds.Dy())).Send()
		return frame
	}
	squareRect := image.Rect(0, 0, SquareSize, SquareSize)
	bordersImage := image.NewRGBA(oldBounds)

	leftOffset := (oldBounds.Dx() - squareRect.Dx()) / 2
	topOffset := (oldBounds.Dy() - squareRect.Dy()) / 2
	rightOffset := oldBounds.Dx() - leftOffset
	bottomOffset := oldBounds.Dy() - topOffset

	// Рисуем квадрат на картинке (область вырезания)
	for x := 0; x < oldBounds.Dx(); x++ {
		for y := 0; y < oldBounds.Dy(); y++ {

			if (x == leftOffset || x == rightOffset) && y >= topOffset && y <= bottomOffset {
				bordersImage.Set(x, y, color.RGBA{R: 180, G: 0, B: 0, A: 255})
			} else if (y == topOffset || y == bottomOffset) && x >= leftOffset && x <= rightOffset {
				bordersImage.Set(x, y, color.RGBA{R: 180, G: 0, B: 0, A: 255})
			} else {
				bordersImage.Set(x, y, frame.At(x, y))
			}
		}
	}

	return bordersImage
}

func (l *Lightmaster) makeEntropyStruct(value uint64) model.EntropyValue {
	return model.EntropyValue{
		Uint64:  value,
		Float64: float64(value) / float64(math.MaxUint64),
		Binary:  fmt.Sprintf("%064b", value),
	}
}
