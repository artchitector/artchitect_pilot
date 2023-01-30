package handler

import (
	"bytes"
	"github.com/artchitector/artchitect/gate/resizer"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"image"
	"image/jpeg"
	"net/http"
)

type CardRequest struct {
	ID uint `uri:"id" binding:"required,numeric"`
}

type ImageRequest struct {
	ID   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"` // size f - full, size m - 2-times smaller dimensions, size s - 4-times smaller dimensions
}

type CardHandler struct {
	cardsRepository cardsRepository
	cache           cache
}

func NewCardHandler(cardsRepository cardsRepository, cache cache) *CardHandler {
	return &CardHandler{cardsRepository, cache}
}

func (lh *CardHandler) Handle(c *gin.Context) {
	var request CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := lh.cache.GetCard(c, uint(request.ID))
	if err != nil {
		log.Error().Err(err).Msgf("[card_handler:Handle] failed to get card(id=%d) from cache", card.ID)
	} else {
		c.JSON(http.StatusOK, card)
		return
	}

	card, found, err := lh.cardsRepository.GetCard(c, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (ch *CardHandler) HandleImage(c *gin.Context) {
	var request ImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cached, err := ch.cache.GetImage(c, uint(request.ID), request.Size)
	if err != nil {
		log.Error().Err(err).Msgf("[card_controller:HandleImage] failed to get cached image")
	} else {
		var newImg []byte
		newImg, err = ch.watermark(cached)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "image/jpeg", newImg)
		return
	}

	log.Fatal().Msgf("STOP")

	card, found, err := ch.cardsRepository.GetCardWithImage(c, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	img, err := resizer.Resize(card.Image.Data, request.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "image/jpeg", img)
}

func (ch *CardHandler) watermark(imgB []byte) ([]byte, error) {
	text := "2000"
	r := bytes.NewReader(imgB)
	img, err := jpeg.Decode(r)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[watermark] failed to decode jpeg")
	}
	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, img.Bounds(), img, image.Point{}, draw.Src)

	fontFace, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[watermark] failed get face font")
	}
	fontDrawer := font.Drawer{
		Dst: rgbaImg,
		Src: image.White,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    128.0,
			Hinting: font.HintingFull,
		}),
	}
	// calculate
	textBounds, _ := fontDrawer.BoundString(text)
	xPosition := (fixed.I(rgbaImg.Rect.Max.X) - fontDrawer.MeasureString(text)) / 2
	log.Info().Msgf(
		"[watermark] maxX: %d, measure: %s, result: %s",
		rgbaImg.Rect.Max.X,
		fontDrawer.MeasureString(text).String(),
		xPosition.String(),
	)
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	yPosition := fixed.I((rgbaImg.Rect.Max.Y)-textHeight.Ceil())/2 + fixed.I(textHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: xPosition,
		Y: yPosition,
	}

	fontDrawer.DrawString(text)
	log.Info().Msgf("[watermark] draw success")
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, rgbaImg, nil); err != nil {
		return []byte{}, errors.Wrapf(err, "[watermark] failed to encode jpeg")
	}

	return buf.Bytes(), nil
}
