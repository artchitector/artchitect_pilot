package entropy

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/soul/resources"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"image"
	"image/color"
	"math"
	"math/bits"
	"os"
	"sync"
	"time"
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
	webcam        *resources.Webcam
	gatekeeper    *Gatekeeper
	lastNFrames   []image.Image
	tags          []string
	selectedWords map[string]int
	counter       int

	entropyMutex         sync.Mutex
	lastEntropyValueUsed bool
	lastEntropyValue     float64
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
	if err := l.gatekeeper.NotifyEntropyPhase(ctx, l.addBordersOnFrame(newFrame), PhaseSource, 0); err != nil {
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
		if _, entropyF, err := l.pipelineEntropy(ctx); err != nil {
			log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoise)
		} else {
			l.entropyMutex.Lock()
			defer l.entropyMutex.Unlock()

			l.lastEntropyValue = entropyF
			l.lastEntropyValueUsed = false
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

func (l *Lightmaster) pipelineEntropy(ctx context.Context) (uint64, float64, error) {
	noiseImage, err := l.sourceToNoise()

	if err != nil {
		return 0, 0.0, errors.Wrapf(err, "[lightmaster] failed to transform source to noise")
	} else if err := l.gatekeeper.NotifyEntropyPhase(ctx, noiseImage, PhaseNoise, 0); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoise)
	}

	shrinkedNoise, err := l.shrinkNoise(noiseImage)
	if err != nil {
		return 0, 0.0, errors.Wrapf(err, "[lightmaster] failed to shrink noise")
	} else if err := l.gatekeeper.NotifyEntropyPhase(ctx, shrinkedNoise, PhaseNoiseShrink, 0); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseNoiseShrink)
	}

	bytesImage, entropyI, err := l.noiseToBytes(shrinkedNoise)
	if err != nil {
		return 0, 0.0, errors.Wrapf(err, "[lightmaster] failed noise2bytes")
	} else if err := l.gatekeeper.NotifyEntropyPhase(ctx, bytesImage, PhaseBytes, entropyI); err != nil {
		log.Error().Err(err).Msgf("[lightmaster] failed to notify gate with phase %s", PhaseBytes)
	}

	entropyF := float64(entropyI) / float64(math.MaxUint64)
	return entropyI, entropyF, nil
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

func (l *Lightmaster) shrinkNoise(noiseImage image.Image) (image.Image, error) {
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

	scale := maxPower - minPower
	for x := 0; x < resultBounds.Dx(); x++ {
		for y := 0; y < resultBounds.Dy(); y++ {
			powerOfPixel := powers[x][y] - minPower // clear extra power

			redPower := math.Round(float64(powerOfPixel) / float64(scale) * 255.0)

			resultImg.SetRGBA(x, y, color.RGBA{R: uint8(redPower), G: uint8(0), B: uint8(0), A: 255})
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
			clr := noiseImage.At(x, y)
			power := clr.(color.RGBA).R
			power = bits.Reverse8(power)

			if power > 128 { // set white
				bytesImage.Set(x, y, color.RGBA{R: power, G: 0, B: 0, A: 255})

				byteIndex := x*ResultSize + y
				entropyAnswer = entropyAnswer | 1<<(63-byteIndex)
			} else {
				bytesImage.Set(x, y, color.RGBA{R: power, G: 0, B: 0, A: 255})
			}
		}
	}

	return bytesImage, entropyAnswer, nil
}

func (l *Lightmaster) testEntropy(entropy uint64) {
	entropyFl := float64(entropy) / float64(math.MaxUint64)

	if l.tags == nil {
		yamlFile, err := os.ReadFile("files/tags_v12.yaml")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		tags := []string{}
		err = yaml.Unmarshal(yamlFile, &tags)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		l.tags = tags
	}

	targetIndex := int(math.Floor(float64(len(l.tags)) * entropyFl))
	selectedTag := l.tags[targetIndex]

	if _, ok := l.selectedWords[selectedTag]; !ok {
		l.selectedWords[selectedTag] = 1
	} else {
		l.selectedWords[selectedTag] += 1
	}
	l.counter += 1

	if l.counter%100 == 0 {
		l.saveWords()
	}
}

func (l *Lightmaster) saveWords() {
	f, err := os.Create("test_result.txt")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("Total words taken = %d\n", l.counter))
	for _, word := range l.tags {
		if counter, found := l.selectedWords[word]; found {
			f.WriteString(fmt.Sprintf("%04d\t%s\n", counter, word))
		} else {
			f.WriteString(fmt.Sprintf("%04d\t%s\n", 0, word))
		}
	}

	log.Info().Msgf("file saved")
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
