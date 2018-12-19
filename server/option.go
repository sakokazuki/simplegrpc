package server

import (
	"github.com/sakokazuki/simplegrpc/config"
	"github.com/sakokazuki/simplegrpc/pubsub"
)

type Option struct {
	Config   config.Config
	PubSuber pubsub.PubSuber
}
