package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bibaroc/launch/watchdir"
)

type configuration struct {
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

func main() {
	buffer, err := ioutil.ReadFile("launch.config.json")
	if err != nil {
		log.Println(err)
	}
	value := configuration{}
	json.Unmarshal(buffer, &value)
	fmt.Println(value)

	cha := make(chan watchdir.ModificationEvent, 10)
	for i := 0; i < 1; i++ {
		go watchdir.StartWatching(i, ".", cha)
	}

	for v := range cha {
		fmt.Println(v)
	}
}
