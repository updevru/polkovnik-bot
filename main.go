package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.Handle("/api/team", API.TeamList()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/users", API.UserList()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/users/{userId}", API.UserGet()).Methods(http.MethodGet)
	router.Handle("/api/team/{teamId}/users/{userId}", API.UserEdit()).Methods(http.MethodPatch)
	router.Handle("/api/team/{teamId}/users/{userId}", API.UserDelete()).Methods(http.MethodDelete)
	router.Handle("/api/team/{teamId}/users", API.UserAdd()).Methods(http.MethodPost)
	router.Handle("/api/team/{teamId}/tasks", API.TaskList()).Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))

	folder, _ := os.Getwd()
	fs := http.FileServer(http.Dir(folder + "\\ui\\build"))
	router.PathPrefix("/").Handler(fs)

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Println("Starting server at", port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func main() {
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

	if false {
		go func() {
			for tick := range ticker.C {
				now := tick.In(time.Local)
				for _, team := range config.Teams {
					log.Info("Process team", team.Title)
					err := processor.ProcessTeamTasks(team, now)
					if err != nil {
						log.Error("Task error", err)
					}
				}
			}
		}()
	}

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
