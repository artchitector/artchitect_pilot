package handler

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type UnityRequest struct {
	Mask string `uri:"mask" binding:"required"`
}

type Response struct {
	Type    string
	Parent  model.Unity
	Unities []model.Unity
	Cards   []model.Card
}

type UnityHandler struct {
	unityRepository unityRepository
	cardsRepository cardsRepository
}

func NewUnityHandler(unityRepository unityRepository, cardsRepository cardsRepository) *UnityHandler {
	return &UnityHandler{unityRepository, cardsRepository}
}

func (uh *UnityHandler) HandleList(c *gin.Context) {
	unities, err := uh.unityRepository.GetRootUnities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unities)
}

func (uh *UnityHandler) HandleUnity(c *gin.Context) {
	var request UnityRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	unity, err := uh.unityRepository.GetUnity(request.Mask)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if unity.Rank == model.Rank100 {
		// get cards
		cards, err := uh.getCards(c, unity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := Response{
			Type:    "cards",
			Parent:  unity,
			Unities: nil,
			Cards:   cards,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// else model.Rank1000 or model.Rank10000 - get subunities
	unities, err := uh.unityRepository.GetChildUnities(unity.Mask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := Response{
		Type:    "unity",
		Parent:  unity,
		Unities: unities,
		Cards:   nil,
	}
	c.JSON(http.StatusOK, response)
}

func (uh *UnityHandler) getCards(ctx context.Context, unity model.Unity) ([]model.Card, error) {
	start, err := strconv.Atoi(strings.ReplaceAll(unity.Mask, "X", "0"))
	if err != nil {
		return []model.Card{}, errors.Wrapf(err, "[unity_repo] failed to get start of mask %s", unity.Mask)
	}
	end, err := strconv.Atoi(strings.ReplaceAll(unity.Mask, "X", "9"))
	if err != nil {
		return []model.Card{}, errors.Wrapf(err, "[unity_repo] failed to get end of mask %s", unity.Mask)
	}
	log.Info().Msgf("[unity_repo] mask %s become range %d-%d", unity.Mask, start, end)
	return uh.cardsRepository.GetCardsByRange(uint(start), uint(end))
}
