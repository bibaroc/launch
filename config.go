package main

import (
	"fmt"
)

//Configuration rappresents basic configuration
type Configuration struct {
	Target []watchedDir  `json:"monitor"`
	Timed  []timedAction `json:"repeat"`
}

//WatchedDir is a watched dir
type watchedDir struct {
	Path      string `json:"watch"`
	Timeout   string `json:"after"`
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

func (c watchedDir) String() string {
	return fmt.Sprintf("WatchedDir=[Path=%v, Timeout=%v, MatchRule=%v, Action=%v]", c.Path, c.Timeout, c.MatchRule, c.Action)
}

func (t timedAction) String() string {
	return fmt.Sprintf("timedAction[Action=%v, Timespan=%v]", t.Action, t.Timespan)
}
