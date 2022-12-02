package infrastructure

import (
	"github.com/rs/zerolog"
)

/*
Cloud helps pr ay and receive gifts by other services
*/
type Cloud struct {
	logger zerolog.Logger
}

// TODO Need to think about cloud-architecture. Now I don't need complex postgres notify/listen.
// I have only one instance of service, and golang channels will be more efficient and simple, than pub/sub architectures.
