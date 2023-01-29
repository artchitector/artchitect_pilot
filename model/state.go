package model

type CreationState struct {
	Version              string
	Seed                 uint
	TagsCount            uint
	Tags                 []string
	LastCardPaintTime    uint // seconds
	CurrentCardPaintTime uint // seconds
	CardID               uint
	EnjoyTime            uint
	CurrentEnjoyTime     uint
}

type LotteryState struct {
	Lottery          Lottery
	EnjoyTotalTime   uint
	EnjoyCurrentTime uint
}

type PrayState struct {
	Queue   uint
	Started bool
}
