package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const usageMessage = "" +
	`
Usage of the Visualization:
	vis [flags(see below)] gcn3 [your_database_name]
	Example: ./server -http=:8080 gcn3 bkeghtbf4tduanale8u0

Flags
	-http=addr: HTTP service address
	Example: ... -http=:8080 ...
	`

var (
	httpFlag = flag.String("http",
		"0.0.0.0:3001",
		"HTTP service address (e.g., ':6060')")
	model     string
	username  string
	password  string
	ipAddress string
	port      int
	dbName    string
	db        *sql.DB
)

func main() {
	parseArgs()
	startServer()
}

func parseArgs() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(2)
	}

	flag.Parse()

	switch flag.NArg() {
	case 2:
		model = flag.Arg(0)
		dbName = flag.Arg(1)
	default:
		flag.Usage()
	}
}

func startServer() {
	connectToDB()
	// summarizeComponentInfo()
	startAPIServer()
}

func connectToDB() {
	getCredentials()

	fmt.Printf("Open trace in db %s\n", dbName)

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		username, password, ipAddress, port, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func getCredentials() {
	var err error

	username = os.Getenv("AKITA_TRACE_USERNAME")
	if username == "" {
		panic(`trace username is not set, use environment variable AKITA_TRACE_USERNAME to set it.`)
	}

	password = os.Getenv("AKITA_TRACE_PASSWORD")
	ipAddress = os.Getenv("AKITA_TRACE_IP")
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}

	portString := os.Getenv("AKITA_TRACE_PORT")
	if portString == "" {
		portString = "3306"
	}
	port, err = strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}
}

func startAPIServer() {
	fs := http.FileServer(http.Dir("../" + model))

	http.Handle("/", fs)
	http.HandleFunc("/dashboard",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../"+model+"/index.html")
		})
	http.HandleFunc("/component",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../"+model+"/index.html")
		})
	http.HandleFunc("/task",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../"+model+"/index.html")
		})

	// http.HandleFunc("/api/task", httpTask)
	http.HandleFunc("/api/trace", httpTrace)
	http.HandleFunc("/api/compnames", httpComponentNames)
	http.HandleFunc("/api/compinfo", httpComponentInfo)
	// http.HandleFunc("/api/ips", handleIPS)w

	fmt.Printf("Listening %s\n", *httpFlag)
	err := http.ListenAndServe(*httpFlag, nil)
	dieOnErr(err)
}

func dieOnErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
