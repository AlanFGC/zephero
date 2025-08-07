package core

type PlayerEvent struct {
	PlayerId  string    `json:"playerId"`
	GameEvent GameEvent `json:"gameEvent"`
}

type GameEvent struct {
	EventId string `json:"eventId"`
	Data    string `json:"data"`
}
