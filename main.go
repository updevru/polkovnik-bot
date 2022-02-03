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
	"polkovnik/job"
	"polkovnik/repository"
	"syscall"
	"time"
)

//go:generate swagger generate spec -o ./ui/public/doc/api/swagger.json

var stdout *bool
var configFile *string
var dbFile *string
var httpPort *string

//go:embed templates
var templates embed.FS

//go:embed ui/build
var UIFiles embed.FS

func init() {

	stdout = flag.Bool("o", false, "Send logs to stdout")
	configFile = flag.String("c", "var/config.json", "Config file")
	dbFile = flag.String("db", "var/data.db", "Database file")
	httpPort = flag.String("p", "8080", "HTTP port for UI")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{})

	if *stdout == true {
		log.SetOutput(os.Stdout)
	}
}

func runWebServer(port string, repository *repository.Repository, history *repository.HistoryRepository, processor *job.Processor) {
	API := api.NewApiHandler(repository, history, processor)

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
	dataRepository := repository.NewRepository(configStorage)
	err := dataRepository.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Database file: ", *dbFile)
	historyStorage, err := repository.CreateHistoryRepository(*dbFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer historyStorage.Close()

	err = app.Migrate(dataRepository)
	if err != nil {
		log.Fatal(err)
		return
	}

	processor := job.NewProcessor(app.NewTemplateEngine("templates", templates), dataRepository, historyStorage)

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		processor.Stop()
		exit <- true
	}()

	fmt.Println("Run scheduler")
	go processor.StartScheduler()

	fmt.Println("Run worker")
	go processor.StartWorker()

	fmt.Println("Run http server")
	go runWebServer(*httpPort, dataRepository, historyStorage, processor)

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			result, err := dataRepository.Flush()
			if err != nil {
				log.Error(fmt.Sprintf("Save config... Error: %s", err.Error()))
			} else if result {
				log.Info("Save config... Ok")
			}
		}
	}()

	<-exit

	fmt.Print("Save config...")
	_, err = dataRepository.Flush()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK")
	fmt.Println("Buy.")
}
