package model

type CurrentState string

const (
	CurrentStateIdle = CurrentState("idle")
)

type LastPainting struct {
	ID      uint
	Caption string
}

type State struct {
	CurrentState CurrentState
	LastPainting LastPainting
}
