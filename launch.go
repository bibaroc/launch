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

	"github.com/fsnotify/fsnotify"
)

func execute(dir watchedDir, quit *sync.WaitGroup) {
	var (
		mx      sync.Mutex
		path, _ = filepath.Abs(dir.Path)
		args    = strings.Split(dir.Action, " ")
		cmd     *exec.Cmd
		exqt    = func(what string, how ...string) func() {
			return func() {
				cmd = exec.Command(what, how...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				err := cmd.Start()
				if err != nil {
					log.Panic("execute error starting:" + dir.Action + err.Error())
				}
			}
		}
		t0, _        = time.ParseDuration("1us")
		timer        = time.AfterFunc(t0, exqt(args[0], args[1:]...))
		watcher, err = fsnotify.NewWatcher()
	)

	defer quit.Done()
	defer watcher.Close()

	if err != nil {
		log.Println("execute:" + err.Error())
		return
	}

	watcher.Add(path)

	for {
		select {
		case event := <-watcher.Events:
			var (
				pwd, _  = os.Getwd()
				t1, err = time.ParseDuration(dir.Timeout)
				m       bool
			)
			if err != nil {
				log.Println("Error parsing time " + err.Error())
				return
			}
			if m, err = regexp.MatchString(dir.MatchRule, strings.TrimLeft(event.Name, pwd)); err != nil {
				log.Println("Error while testing modification event" + err.Error())
				return
			}
			if m {
				log.Println(event.Name, "matched", dir.MatchRule)
				go func() {
					mx.Lock()
					defer mx.Unlock()
					timer.Stop()
					timer = time.AfterFunc(t1, func() {

						if err := cmd.Process.Kill(); err != nil {
							if !os.IsPermission(err) {
								log.Println("Error killing application" + err.Error())
							}
						}

						exqt(args[0], args[1:]...)()
					})
				}()
			}
		case err := <-watcher.Errors:
			log.Println("Error: ", err)
		}
	}
}

func main() {

	var (
		command     = flag.String("exec", "", "The comand you want to be launched and repeated on the current directoy everytime i changes")
		application = Configuration{}
	)

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

	var wg sync.WaitGroup
	log.Println(application)
	for _, v := range application.Target {
		wg.Add(1)
		go execute(v, &wg)
	}
	wg.Wait()
}
