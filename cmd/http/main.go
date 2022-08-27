package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/farrej10/srtl/configs"
	"github.com/farrej10/srtl/internal/shortener"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger
var short shortener.IShortener
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
	short, err = shortener.NewShortener(shortener.Config{
		Logger: *sugar,
		Host:   host,
		Port:   port,
		Home:   configs.Https + "www." + host + "/"})
	if err != nil {
		panic(err)
	}
}

func main() {
	sugar.Info("Starting Shortener")
	http.HandleFunc("/", short.ShortenLink)
	http.ListenAndServe(":"+port, nil)
}
