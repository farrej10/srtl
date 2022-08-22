package main

import (
	"net/http"

	"github.com/farrej10/srtl/internal/shortener"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
}

func main() {
	sugar.Info("Starting Shortener")
	s := shortener.NewShortener(shortener.Config{Logger: *sugar})
	http.HandleFunc("/", s.ShortenLink)
	http.ListenAndServe(":8080", nil)
}
