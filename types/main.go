package types
const PLAY = "play"
const STOP = "stop"

type Rule struct {
	Interval int
	Days     []string
	From    int
	To      int
}

