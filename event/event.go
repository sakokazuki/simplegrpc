package event

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type MetaData struct {
	Type string `json:"type"`
}

//要はMeta.Typeがイベント名でMeta.Dataがメッセージ
type Payload struct {
	Meta MetaData        `json:"meta"`
	Data json.RawMessage `json:"data"`
}

func (pl Payload) String() string {
	s, err := json.Marshal(&pl.Data)
	if err != nil {
		log.Error().Err(err).Msg("faied to marshal payload data")
	}
	return "meta: " + pl.Meta.Type + "  data: " + string(s)
}
