package model

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const (
	UnityStateEmpty         = "empty"
	UnityStateSkipped       = "skipped"
	UnityStateUnified       = "unified"
	UnityStateReunification = "reunification"
)

/*
Понятие "Единство/Unity". Объединяет в себя какую-то группу карточек и имеет своё общее изображение.
Характеризуется маской единства, например:
- десятитысяча - Rank10000. Маска 1**** (это объединит карточки с 10000 до 19999)
- тысяча - Rank1000. Маска 12*** (это объединит карточки с 12000 до 12999)
- сотня - Rank100. Маска 123** (это объединит карточки с 12300 до 12399)
*/
type Unity struct {
	Mask      string `gorm:"primaryKey"`
	Rank      uint
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
	Leads     string  // leads is json array [56123, 56690, ...]. leads used for unity combination in single picture
	Children  []Unity `gorm:"-" json:"-"`
}

func (u *Unity) String() string {
	return fmt.Sprintf("%s/%s", u.Mask, u.State)
}

func (u *Unity) Start() uint {
	num, err := strconv.Atoi(strings.ReplaceAll(u.Mask, "X", "0"))
	if err != nil {
		log.Error().Err(err).Send()
		return 0
	}
	return uint(num)
}
