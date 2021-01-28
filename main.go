package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"teamBot/adapter/storage"
	"teamBot/app"
	"teamBot/job"
	"time"
)

var stdout *bool
var configFile *string

func init() {

	stdout = flag.Bool("o", false, "Send logs to stdout")
	configFile = flag.String("c", "var/config.json", "Config file")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{})

	if *stdout == true {
		log.SetOutput(os.Stdout)
	}
}

func main() {
	configStorage := storage.NewConfigFile(*configFile)
	config, err := configStorage.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	processor := job.Processor{
		Tpl: app.NewTemplateEngine("templates"),
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(time.Minute)

	go func() {
		<-signals
		ticker.Stop()
		exit <- true
	}()

	go func() {
		for tick := range ticker.C {
			now := tick.In(time.Local)
			for _, team := range config.Teams {
				log.Info("Process team", team.Title)
				err := processor.ProcessTeamTasks(&team, now)
				if err != nil {
					log.Error("Task error", err)
				}
			}
		}
	}()

	fmt.Println("Running...")
	<-exit

	fmt.Print("Save config...")
	err = configStorage.Update(config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK")

	fmt.Println("Buy.")
}
