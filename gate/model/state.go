package model

type LastPainting struct {
	ID uint64
}

type State struct {
	LastPainting LastPainting
}
