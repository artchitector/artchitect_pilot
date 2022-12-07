package model

type Pray interface {
	Name() string
}
type Gift interface {
	Name() string
	Error() error
}

type PrayCallback func(gift Gift) error
type GiftCallback func(pray Pray) (Gift, error)
