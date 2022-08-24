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
	sugar.Info("Starting Shortener")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" || host == "" {
		panic("port or host variables not found")
	}
	s, err := shortener.NewShortener(shortener.Config{Logger: *sugar, Host: host, Port: port})
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/l/", s.ShortenLink)
	http.ListenAndServe(":"+port, nil)
}
