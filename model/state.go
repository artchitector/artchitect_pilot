package model

type CreationState struct {
	Version              string
	Seed                 uint
	TagsCount            uint
	Tags                 []string
	LastCardPaintTime    uint // seconds
	CurrentCardPaintTime uint // seconds
	CardID               uint
}
