package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zaunerc/cntrbrowserd/consul"

	auth "github.com/abbot/go-http-auth"
	"github.com/urfave/cli"
)

var consulUrl string

func init() {
	// By default logger is set to write to stderr device.
	//log.SetOutput(os.Stdout)
}

func main() {

	var httpPort int
	var htpasswd string

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "httpPort, p",
			Value:       80,
			Usage:       "Listen on port `PORT` for HTTP connections",
			Destination: &httpPort,
		},
		cli.StringFlag{
			Name: "htpasswd, s",
			//Value:       "/usr/local/etc/cntrinfod/htpasswd",
			Usage:       "Use htpasswd `FILE`. Enables user authentication.",
			Destination: &htpasswd,
		},
		cli.StringFlag{
			Name:        "consulUrl, c",
			Usage:       "Register container with consul at `URL`. Enables consul registration.",
			Destination: &consulUrl,
		},
	}

	app.Email = "christoph.zauner@NLLK.net"
	app.Author = "Christoph Zauner"
	app.Version = "0.1.0-alpha"
	// cntrinfod, cntinfod
	app.Usage = "Container Browser Daemon: Container-centric HTTP frontend for consul"

	app.Action = func(c *cli.Context) error {

		consul.ScheduleCleanUpTask(consulUrl)

		fmt.Printf("Starting HTTP daemon on port %d...\n", httpPort)
		http.HandleFunc("/", protect(handler, htpasswd))
		http.ListenAndServe(":"+strconv.Itoa(httpPort), nil)

		return nil
	}

	app.Run(os.Args)
}

func protect(handlerFunc http.HandlerFunc, htpasswdPath string) http.HandlerFunc {

	if htpasswdPath != "" {
		// read from htpasswd file
		htpasswd := auth.HtpasswdFileProvider(htpasswdPath)
		authenticator := auth.NewBasicAuthenticator("cntrinfod htpasswd realm", htpasswd)
		return auth.JustCheck(authenticator, handlerFunc)
	} else {
		return func(w http.ResponseWriter, r *http.Request) {
			handlerFunc(w, r)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	t, error := template.ParseFiles("index.html")

	currentDateAndTime := time.Now().Format("2006-01-02 15:04:05")
	containers := consul.FetchContainerData(consulUrl)

	vars := map[string]interface{}{
		"CurrentDateAndTime": currentDateAndTime,
		"Containers":         containers,
	}

	error = t.Execute(w, vars)

	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
	}
}
