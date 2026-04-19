package ws

// Channel is similar to a `room|group|realm` for communication.
// e.g. `lobby`, `game.{game_id}`, `gametv.{game_id}`, `gametv.{game_id}.chat` etc.
type Channel string

func (ch Channel) String() string {
	return string(ch)
}

func channelSafePrint(ch *Channel) string {
	if ch == nil {
		return ""
	}
	return string(*ch)
}
