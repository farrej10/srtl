package main

import (
	"math/rand"
	"net/http"
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
	sugar.Info("Starting Shortener")
	s, err := shortener.NewShortener(shortener.Config{Logger: *sugar})
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/l/", s.ShortenLink)
	http.ListenAndServe(":8080", nil)
}
