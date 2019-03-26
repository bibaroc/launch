package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/bibaroc/launch/config"
	"github.com/bibaroc/launch/watchdir"
)

type processes struct {
	processes map[int]exec.Cmd
}

func main() {
	//var declaration
	var (
		command      = flag.String("exec", "", "The comand you want to be launched and repeated on the current directoy everytime i changes")
		application  = config.Configuration{}
		eventChannel = make(chan watchdir.ModificationEvent, 10)
	)
	//init
	flag.Parse()
	if *command != "" {
		application.Target = append(application.Target, config.WatchedDir{Path: ".", Timeout: 300, MatchRule: "*", Action: *command})
	} else {
		path, err := filepath.Abs("launch.config.json")
		if err != nil {
			log.Fatalln("init failed: ", err)
		}
		buffer, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalln("init failed: ", err)
		}
		json.Unmarshal(buffer, &application)
	}

	for i := 0; i < 1; i++ {
		go watchdir.StartWatching(i, ".", eventChannel)
	}

	for v := range eventChannel {
		fmt.Println(v)
	}
}
