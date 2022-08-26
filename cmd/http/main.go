package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/farrej10/srtl/internal/shortener"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	sugar.Infof("Starting Shortener on Host: %s", host)
	if port == "" || host == "" {
		panic("port or host variables not found")
	}
	s, err := shortener.NewShortener(shortener.Config{Logger: *sugar, Host: host, Port: port})
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", s.ShortenLink)
	http.ListenAndServe(":"+port, nil)
}
