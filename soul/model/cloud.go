package model

type Pray struct {
	Name    string
	Payload interface{}
}
type Gift struct {
	Name    string
	Payload interface{}
	Error   error
}

type PrayCallback func(gift Gift) error
type GiftCallback func(pray Pray) (Gift, error)
