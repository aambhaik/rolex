package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/robfig/cron"
)

var runRolex = true

var (
	timerurl = flag.String("timerurl", "localhost:9980", "URL for the timer service")
)

func main() {
	flag.Parse()

	f, err := os.OpenFile("/Users/aambhaik/log/rolex.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	c := cron.New()

	c.AddFunc("@every 10s", func() {
		if runRolex {
			timecheck(*timerurl)
		}
	})

	c.Start()

	httpRoute()
	log.Printf("Rolex service started successfully")
}

func httpRoute() {
	router := httprouter.New()
	router.GET("/ping", PingHandler)
	router.GET("/start", StartHandler)

	log.Printf("Starting Rolex service on port", 9985)
	http.ListenAndServe(":9985", router)
}

func PingHandler(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writer.WriteHeader(200)
}

func StartHandler(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	log.Printf("Re-starting timer service polling from rolex")
	runRolex = true
}

func timecheck(timerurl string) {
	resp, err := http.Get(timerurl)
	if err != nil {
		// handle err
		log.Printf("The timer service is not available")
		runRolex = false
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp != nil {
		if resp.StatusCode == http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			log.Printf("The rolex time now is %v", string(bodyBytes))
		}
	}
}
