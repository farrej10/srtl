package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/farrej10/srtl/configs"
	pebbleadapter "github.com/farrej10/srtl/internal/adapters/pebble_adapter"
	"github.com/farrej10/srtl/internal/ports"
	"github.com/farrej10/srtl/internal/shortener"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger
var short shortener.IShortener
var db ports.IDatabaseAccessor
var port string

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	rand.Seed(time.Now().UnixNano())
	port = os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" || host == "" {
		panic("port or host variables not found")
	}
	var err error
	db, err = pebbleadapter.NewPebbleDb("./dbPebble", sugar)
	if err != nil {
		panic(err)
	}
	short, err = shortener.NewShortener(shortener.Config{
		Logger: sugar,
		Host:   host,
		Port:   port,
		Home:   configs.Https + "www." + host + "/",
		Db:     db,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	sugar.Info("Starting Shortener")
	http.HandleFunc("/", short.ShortenLink)
	http.ListenAndServe(":"+port, nil)
	defer db.Close()
}
