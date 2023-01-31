package watermark

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Watermark struct {
	font   *truetype.Font
	catImg image.Image
}

func NewWatermark() *Watermark {
	return &Watermark{}
}

// RU: Так как тут сложно, многие комментарии будут на русском

// AddWatermark - adds watermark over the original image (with card number and cat icon)
func (w *Watermark) AddWatermark(originalImage image.Image, cardID uint) (image.Image, error) {
	// RU: прогружаем ресурсы (шрифт + иконка кота), если они еще не загружены
	if err := w.loadResources(); err != nil {
		return nil, errors.Wrap(err, "[watermark] failed to load resources")
	}
	// RU: Формируем новый RGBA-холст и копируем начальную картину на него
	finalImage := image.NewRGBA(originalImage.Bounds())
	draw.Draw(finalImage, originalImage.Bounds(), originalImage, image.Point{}, draw.Src)

	// RU: Вотермарка в виде image.Image, накладывается на холст в правый нижний угол с отступом
	watermarkImg := w.makeWatermarkImage(originalImage.Bounds(), cardID)
	padding := fixed.I(finalImage.Bounds().Max.X / 30) // RU: Отступ = 1/30 от ширины холста
	rightBottomPoint := image.Point{
		// RU: Тут всё было подобрано наугад. Почему отрицательные - не знаю
		X: -1 * (fixed.I(finalImage.Bounds().Max.X) - padding - fixed.I(watermarkImg.Bounds().Dx())).Ceil(),
		Y: -1 * (fixed.I(finalImage.Bounds().Max.Y) - padding - fixed.I(watermarkImg.Bounds().Dy())).Ceil(),
	}
	draw.Draw(finalImage, finalImage.Bounds(), watermarkImg, rightBottomPoint, draw.Over)
	return finalImage, nil
}

// makeWatermarkImage - prepare watermark image (with transparent black background, number and cat)
// we need bounds to select font-size for current image (there are old low-res images and new hi-res)
func (w *Watermark) makeWatermarkImage(bounds image.Rectangle, cardID uint) image.Image {
	var size float64
	if bounds.Dx() < 1000 {
		// very old
		size = 30.0
	} else if bounds.Dx() < 2000 {
		// old low-res (1024x1536)
		size = 43.0
	} else {
		// new hi-res
		size = 86.0
	}

	// RU: Это рисователь шрифта на картинке
	fontDrawer := font.Drawer{
		Src: image.NewUniform(color.RGBA{200, 200, 200, 255}), // gray font color
		Face: truetype.NewFace(w.font, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
		}),
	}

	text := fmt.Sprintf("#%d", cardID)            // RU: Все номера карточек начинаются с #, там исторически пошло
	textBounds, _ := fontDrawer.BoundString(text) // RU: Тут объект, из которого можно получить итоговые размеры надписи в пикселях
	textWidth := fontDrawer.MeasureString(text)
	textHeight := textBounds.Max.Y - textBounds.Min.Y // Это же является высотой и шириной кота

	// RU: Подготовили изображение кота, который в левой части вотермарки
	catImg := w.prepareCatForWatermark(textHeight.Ceil())
	// RU: Готовим холст самой вотермарки. Берём ширину текст + ширину кота (умноженную на 7/6, чтобы был отступ)
	watermarkImg := image.NewRGBA(image.Rect(
		0,
		0,
		(textWidth + textHeight*7/6).Ceil(),
		textHeight.Ceil(),
	))
	blackTransparent := color.RGBA{0, 0, 0, 128} // RU: Фон вотермарки
	// RU: Заполняем холст вотермарки фоном
	draw.Draw(watermarkImg, watermarkImg.Bounds(), &image.Uniform{blackTransparent}, image.Point{}, draw.Src)
	// RU: Рисуем сам текст на холсте
	fontDrawer.Dst = watermarkImg
	fontDrawer.Dot = fixed.Point26_6{
		X: textHeight * 7 / 6,                   // RU:Отступаем на расстояние кота и его отступа
		Y: fixed.I(watermarkImg.Bounds().Max.Y), // RU:Рисуется от низа вотермарки вверх
	}
	fontDrawer.DrawString(text)

	// RU: Теперь добавляем кота на холст в левую часть (он по высоте как и холст)
	draw.Draw(watermarkImg, watermarkImg.Bounds(), catImg, image.Point{X: 0, Y: 0}, draw.Over)
	return watermarkImg
}

func (w *Watermark) prepareCatForWatermark(watermarkHeight int) image.Image {
	catImgResized := image.NewRGBA(image.Rect(0, 0, watermarkHeight, watermarkHeight))
	draw.NearestNeighbor.Scale(catImgResized, catImgResized.Rect, w.catImg, w.catImg.Bounds(), draw.Over, nil)
	return catImgResized
}

func (w *Watermark) loadResources() error {
	if w.font == nil {
		log.Info().Msgf("[watermark] load resources: font")
		// load font from file and parse it
		fontFile := "./files/conso.ttf"
		fontData, err := os.ReadFile(fontFile)
		if err != nil {
			return errors.Wrapf(err, "[watermark] failed to load font from file %s", fontFile)
		}
		fontFace, err := freetype.ParseFont(fontData)
		if err != nil {
			return errors.Wrapf(err, "[watermark] failed to parse font from file %s", fontFile)
		}
		w.font = fontFace
	}
	if w.catImg == nil {
		log.Info().Msgf("[watermark] load resources: cat")
		catImgData, err := os.ReadFile("./files/watermark.png")
		if err != nil {
			return errors.Wrapf(err, "[watermark] failed to load cat image")
		}
		r := bytes.NewReader(catImgData)
		catImg, err := png.Decode(r)
		if err != nil {
			return errors.Wrapf(err, "[watermark] failed to decode cat image")
		}
		w.catImg = catImg
	}
	return nil
}
