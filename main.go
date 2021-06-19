package main

import (
	"embed"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"polkovnik/adapter/storage"
	"polkovnik/api"
	"polkovnik/app"
	"polkovnik/domain"
	"polkovnik/job"
	"polkovnik/repository"
	"syscall"
	"time"
)

var stdout *bool
var configFile *string
var httpPort *string

//go:embed templates
var templates embed.FS

//go:embed ui/build
var UIFiles embed.FS

func init() {

	stdout = flag.Bool("o", false, "Send logs to stdout")
	configFile = flag.String("c", "var/config.json", "Config file")
	httpPort = flag.String("p", "8080", "HTTP port for UI")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{})

	if *stdout == true {
		log.SetOutput(os.Stdout)
	}
}

func runWebServer(port string, config *domain.Config) {
	API := api.NewApiHandler(repository.NewRepository(config))

	server := http.Server{
		Addr:    ":" + port,
		Handler: api.CreateRouter(API, UIFiles),
	}

	fmt.Println("Starting server at", port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Config file: ", *configFile)
	configStorage := storage.NewConfigFile(*configFile)
	config, err := configStorage.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = app.Migrate(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	processor := job.Processor{
		Tpl: app.NewTemplateEngine("templates", templates),
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
				log.Info("Process team ", team.Title)
				err := processor.ProcessTeamTasks(team, now)
				if err != nil {
					log.Error("Task error ", err)
				}
			}
		}
	}()

	fmt.Println("Running...")
	go runWebServer(*httpPort, config)
	<-exit

	fmt.Print("Save config...")
	err = configStorage.Update(config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK")

	fmt.Println("Buy.")
}
