package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bibaroc/launch/watchdir"
)

func execute(dir watchedDir, quit chan string) {
	var (
		events    = make(chan watchdir.ModificationEvent, 10)
		mx        sync.Mutex
		path, err = filepath.Abs(dir.Path)
		args      = strings.Split(dir.Action, " ")
		cmd       *exec.Cmd
		exqt      = func(what string, how ...string) func() {
			return func() {
				cmd = exec.Command(what, how...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				// cmd.Stdin = os.Stdin
				err = cmd.Start()
				if err != nil {
					quit <- "execute error starting:" + dir.Action + err.Error()
				}
			}
		}
		t0, _ = time.ParseDuration("1us")
		timer = time.AfterFunc(t0, exqt(args[0], args[1:]...))
	)
	if err != nil {
		quit <- "execute:" + err.Error()
		return
	}

	go watchdir.StartWatching(path, events)

	for v := range events {
		t1, err := time.ParseDuration(dir.Timeout)
		if err != nil {
			quit <- "Error parsing time " + err.Error()
		}

		log.Println("Detected", v.Action, "of", v.FilePath)
		if m, err := regexp.MatchString(dir.MatchRule, v.FilePath); err != nil {
			quit <- "Error while testing modification event" + err.Error()
		} else {
			if m {
				//If the process is already dead a permission error is raised
				//I think it's because you are trying to access illegal memory
				//If code reaches this point i know i can start a process, there is no reason i couldnt quit itc
				log.Println(v.FilePath, "matched", dir.MatchRule)
				go func() {
					mx.Lock()
					defer mx.Unlock()
					timer.Stop()
					timer = time.AfterFunc(t1, func() {

						if err := cmd.Process.Kill(); err != nil {
							if !os.IsPermission(err) {
								quit <- "Error killing application" + err.Error()
							}
						}

						exqt(args[0], args[1:]...)()
					})
				}()
			}
		}
	}
}

func main() {
	//var declaration
	var (
		command     = flag.String("exec", "", "The comand you want to be launched and repeated on the current directoy everytime i changes")
		application = Configuration{}
	)
	//init
	flag.Parse()
	if *command != "" {
		application.Target = append(application.Target, watchedDir{Path: ".", Timeout: "300ms", MatchRule: ".*", Action: *command})
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
