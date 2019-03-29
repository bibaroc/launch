package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bibaroc/launch/config"
	"github.com/bibaroc/launch/watchdir"
)

func execute(dir config.WatchedDir, quit chan string) {
	var (
		events    = make(chan watchdir.ModificationEvent, 10)
		path, err = filepath.Abs(dir.Path)
		args      = strings.Split(dir.Action, " ")
		cmd       = exec.Command(args[0], args[1:]...)
	)
	if err != nil {
		quit <- "execute:" + err.Error()
		return
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Start()

	if err != nil {
		quit <- "execute error starting:" + dir.Action + err.Error()
	}

	go watchdir.StartWatching(path, events)

	for v := range events {
		log.Println("Detected", v.Action, "of", v.FilePath)
	}
}

func main() {
	//var declaration
	var (
		command     = flag.String("exec", "", "The comand you want to be launched and repeated on the current directoy everytime i changes")
		application = config.Configuration{}
	)
	//init
	flag.Parse()
	if *command != "" {
		application.Target = append(application.Target, config.WatchedDir{Path: ".", Timeout: 300, MatchRule: "*", Action: *command})
	} else {
		path, err := filepath.Abs("launch.config.json")
		if err != nil {
			log.Fatalln("init failed:", err)
		}
		buffer, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalln("init failed:", err)
		}
		json.Unmarshal(buffer, &application)
	}

	log.Println(application)

	errors := make(chan string)
	for _, v := range application.Target {
		go execute(v, errors)
	}

	for v := range errors {
		log.Fatalln(v)
	}

}
