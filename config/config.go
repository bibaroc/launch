package config

import (
	"fmt"
)

//Configuration rappresents basic configuration
type Configuration struct {
	Target []WatchedDir  `json:"monitor"`
	Timed  []timedAction `json:"repeat"`
}

//WatchedDir is a watched dir
type WatchedDir struct {
	Path      string `json:"watch"`
	Timeout   int    `json:"after"`
	MatchRule string `json:"test"`
	Action    string `json:"do"`
}

type timedAction struct {
	Action   string `json:"do"`
	Timespan int    `json:"every"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("Configuration=[Target=%v, Timed=%v]", c.Target, c.Timed)
}

func (c WatchedDir) String() string {
	return fmt.Sprintf("WatchedDir=[Path=%v, Timeout=%v, MatchRule=%v, Action=%v]", c.Path, c.Timeout, c.MatchRule, c.Action)
}

func (t timedAction) String() string {
	return fmt.Sprintf("timedAction[Action=%v, Timespan=%v]", t.Action, t.Timespan)
}
