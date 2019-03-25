package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bibaroc/launch/config"
	"github.com/bibaroc/launch/watchdir"
)

func main() {
	buffer, err := ioutil.ReadFile("launch.config.json")
	if err != nil {
		log.Println(err)
	}
	value := config.Configuration{}
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
