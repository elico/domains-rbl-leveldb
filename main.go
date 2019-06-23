package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	database *leveldb.DB
	logger   *log.Logger
	port     string
	dbFile   string
)

func cacheResponseFor(w http.ResponseWriter, r *http.Request, seconds int) {
	cacheUntil := time.Now().UTC().Add(time.Duration(seconds) * time.Second).Format(http.TimeFormat)
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(seconds))
	w.Header().Set("Expires", cacheUntil)
}

func dontCacheResponse(w http.ResponseWriter, r *http.Request) {
	cacheUntil := time.Now().UTC().Add(time.Duration(-3600) * time.Second).Format(http.TimeFormat)
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	w.Header().Set("Expires", cacheUntil)
}

func init() {
	logger = log.New(os.Stderr, "DRBL: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting up")
	flag.StringVar(&port, "port", "8080", "Port for the web service to listen")
	flag.StringVar(&dbFile, "db", "blacklist.db", "Path for leveldb file")
	flag.Parse()
}

func main() {

	var err error

	database, err = leveldb.OpenFile(dbFile, nil)
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		database.Close()
		os.Exit(1)
	}()
	defer database.Close()

	err = database.Put([]byte("PING"), []byte("PONG"), nil)
	if err != nil {
		panic(err)
	}

	_, err = database.Get([]byte("PING"), nil)
	if err != nil {
		panic(err)
	}

	err = database.Delete([]byte("PING"), nil)
	if err != nil {
		panic(err)
	}

	logger.Println("Database File CONNECTED:", dbFile)

	m := http.NewServeMux()

	// All URLs will be handled by this function
	m.HandleFunc("/insert/", insert)
	m.HandleFunc("/search/", search)
	m.HandleFunc("/recursiveSearch/", recursiveSearch)
	m.HandleFunc("/delete/", delete)

	m.HandleFunc("/whitelist/", delete)
	m.HandleFunc("/blacklist/", insert)
	m.HandleFunc("/test/", search)
	m.HandleFunc("/check/", search)

	m.HandleFunc("/batch/insert/", batchInsert)
	m.HandleFunc("/batch/delete/", batchDelete)

	logger.Println("Starting Web Service... on PORT", port)
	logger.Fatal(http.ListenAndServe(":8080", m))
}
