package shortener

import (
	"net/http"

	"go.uber.org/zap"
)

type IShortener interface {
	ShortenLink(http.ResponseWriter, *http.Request)
}

type (
	shortener struct {
		logger zap.SugaredLogger
	}
	Config struct {
		Logger zap.SugaredLogger
	}
)

func NewShortener(config Config) shortener {
	return shortener{logger: config.Logger}
}

func (s shortener) ShortenLink(http.ResponseWriter, *http.Request) {
	s.logger.Info("Hit API")
}
