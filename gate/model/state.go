package model

type LastPainting struct {
	ID      uint
	Caption string
}

type State struct {
	LastPainting LastPainting
}
