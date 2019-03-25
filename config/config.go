package config

import (
	"fmt"
)

//Configuration rappresents basic configuration
type Configuration struct {
	Target []watchedDir  `json:"monitor"`
	Timed  []timedAction `json:"repeat"`
	Action string        `json:"onStart"`
}

type watchedDir struct {
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
	return fmt.Sprintf("Configuration=[Action=%v, Target=%v, Timed=%v]", c.Action, c.Target, c.Timed)
}

func (c watchedDir) String() string {
	return fmt.Sprintf("watchedDir=[Path=%v, Timeout=%v, MatchRule=%v, Action=%v]", c.Path, c.Timeout, c.MatchRule, c.Action)
}

func (t timedAction) String() string {
	return fmt.Sprintf("timedAction[Action=%v, Timespan=%v]", t.Action, t.Timespan)
}
