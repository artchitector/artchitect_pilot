package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"image/jpeg"
	"net/http"
)

const (
	SizeF  = "f"
	SizeM  = "m"
	SizeS  = "s"
	SizeXS = "xs"
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
}

func NewCardHandler(cardsRepository cardsRepository) *CardHandler {
	return &CardHandler{cardsRepository: cardsRepository}
}

func (lh *CardHandler) Handle(c *gin.Context) {
	var request CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	card, found, err := ch.cardsRepository.GetCard(c, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	img, err := ch.reduceImage(card.Image, request.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "image/jpeg", img)
}

func (ch *CardHandler) reduceImage(rawImg []byte, size string) ([]byte, error) {
	r := bytes.NewReader(rawImg)
	img, err := jpeg.Decode(r)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[card_handler] failed to decode jpeg")
	}

	var height, width uint
	switch size {
	case SizeF:
		// nothing to do, image already full
		return rawImg, nil
	case SizeM:
		width = uint(img.Bounds().Size().X / 2)
		height = uint(img.Bounds().Size().Y / 2)
	case SizeS:
		width = uint(img.Bounds().Size().X / 4)
		height = uint(img.Bounds().Size().Y / 4)
	case SizeXS:
		width = uint(img.Bounds().Size().X / 8)
		height = uint(img.Bounds().Size().Y / 8)
	default:
		// TODO сделать из этого ответ bad-requst, если такое пришло
		return []byte{}, errors.Errorf("[card_hadler] wrong size %s", size)
	}

	img = resize.Resize(width, height, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return []byte{}, errors.Wrapf(err, "[card_handler] failed to encode jpeg")
	}

	return buf.Bytes(), nil
}
