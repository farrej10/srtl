package shortener

import (
	"encoding/json"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/farrej10/srtl/configs"
	"github.com/farrej10/srtl/internal/models"
	"github.com/linxGnu/grocksdb"
	"go.uber.org/zap"
)

type IShortener interface {
	ShortenLink(http.ResponseWriter, *http.Request)
}

type (
	shortener struct {
		logger zap.SugaredLogger
		db     *grocksdb.DB
		tmpl   *template.Template
		host   string
		port   string
		home   string
	}
	Config struct {
		Logger zap.SugaredLogger
		Host   string
		Port   string
		Home   string
	}
)

func NewShortener(config Config) (shortener, error) {
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)

	db, err := grocksdb.OpenDbWithTTL(opts, "./db", 86400)
	if err != nil {
		return shortener{}, err
	}
	tmpl, err := template.ParseFiles("./templates/base.html")
	if err != nil {
		return shortener{}, err
	}
	return shortener{logger: config.Logger, db: db, tmpl: tmpl, host: config.Host, port: config.Port, home: config.Home}, err
}

func (s shortener) ShortenLink(rw http.ResponseWriter, req *http.Request) {
	if req.Method == configs.MethodGet {
		s.redirect(rw, req)
	}
	if req.Method == configs.MethodPost {
		s.createLink(rw, req)
	}
}

// redirects to url if found otherwise return root
func (s shortener) redirect(rw http.ResponseWriter, req *http.Request) {
	key := strings.TrimLeft(req.URL.Path, "/")
	if key == "" {
		http.Redirect(rw, req, s.home, http.StatusFound)
	} else if s.validate(key) {
		val, err := s.db.Get(grocksdb.NewDefaultReadOptions(), []byte(key))
		defer val.Free()
		if err != nil {
			s.logger.Error("error during get from rockdb")
			http.Redirect(rw, req, s.home, http.StatusFound)
		} else if !val.Exists() {
			s.logger.Warnf("key not found %s", key)
			http.Redirect(rw, req, s.home, http.StatusFound)
		} else {
			s.logger.Debugf("key found")
			http.Redirect(rw, req, string(val.Data()), http.StatusFound)
		}
	} else {
		s.logger.Debugf("invalid path %s", req.URL.Path)
		http.Redirect(rw, req, s.home, http.StatusFound)
	}
}

func (s shortener) createLink(rw http.ResponseWriter, req *http.Request) {
	// parse incoming request
	if req.Header.Get(configs.ContentType) == configs.ApplicationJson {
		s.handleFromJson(rw, req)
		return
	} else if req.Header.Get(configs.ContentType) == configs.ApplicationForm {
		s.handleFromForm(rw, req)
		return
	}
}

// validate that the path only contains a-zA-z
func (s shortener) validate(path string) bool {
	s.logger.Debug(path)
	for idx, r := range path {
		// max link length 5
		if idx > 5 {
			return false
		}
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// creates a key of len maxlen make up of the set of runes
func createKey() []byte {
	b := make([]byte, configs.MaxPathLen)
	for i := range b {
		b[i] = configs.ValidRunes[rand.Intn(len(configs.ValidRunes))]
	}
	return b
}

// gets a key that is not already in use
func (s shortener) getKey() ([]byte, error) {
	key := createKey()

	value, err := s.db.Get(grocksdb.NewDefaultReadOptions(), key)
	defer value.Free()
	if err != nil {
		s.logger.Error("error during get from rockdb")
		return nil, err
	}
	// create random key until its not already taken
	for value.Exists() {
		key = createKey()
		value, err = s.db.Get(grocksdb.NewDefaultReadOptions(), key)
		if err != nil {
			s.logger.Error("error during get from rockdb")
			return nil, err
		}
	}
	return key, nil
}

// when getting data from a form it will redirect user to page rather than returning raw data
func (s shortener) handleFromForm(rw http.ResponseWriter, req *http.Request) {
	key, err := s.getKey()
	if err != nil {
		s.logger.Error("failed to get key")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = req.ParseForm()
	if err != nil {
		s.logger.Error("unable to parse form")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	formLink := req.PostForm.Get("link")
	s.logger.Debugf("formLink: %s", formLink)
	if !strings.HasPrefix(formLink, configs.Https) && !strings.HasPrefix(formLink, configs.Http) {
		formLink = configs.Https + formLink
	}
	formLink = strings.Split(formLink, "#")[0]
	incomingUrl, err := url.ParseRequestURI(formLink)
	s.logger.Debugf("incomingUrl: %s", incomingUrl)
	if err != nil {
		s.logger.Error("bad link")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	link := incomingUrl.String()
	s.logger.Debugf("key: %s,value: %s", string(key), link)
	err = s.db.Put(grocksdb.NewDefaultWriteOptions(), key, []byte(link))
	if err != nil {
		s.logger.Error("error during setting rockdb")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	short := s.host + "/" + string(key)
	data := models.ResponseBody{
		Link:  link,
		Short: short,
	}
	s.tmpl.Execute(rw, data)
}

// json queries will return json
func (s shortener) handleFromJson(rw http.ResponseWriter, req *http.Request) {
	key, err := s.getKey()
	if err != nil {
		s.logger.Error("failed to get key")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	incoming := models.CreateLink{}
	err = json.NewDecoder(req.Body).Decode(&incoming)
	defer req.Body.Close()
	if err != nil {
		s.logger.Error("error decoding request")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// drop fragment if sent
	incomingUrl, err := url.ParseRequestURI(strings.Split(incoming.Link, "#")[0])
	// validate uri
	if err != nil {
		s.logger.Error("bad link")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	link := incomingUrl.String()
	s.logger.Debugf("key: %s,value: %s", string(key), link)
	err = s.db.Put(grocksdb.NewDefaultWriteOptions(), key, []byte(link))
	if err != nil {
		s.logger.Error("error during setting rockdb")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	short := configs.Https + s.host + "/" + string(key)
	data := models.ResponseBody{
		Link:  link,
		Short: short,
	}
	rw.Header().Set(configs.ContentType, configs.ApplicationJson)
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(data)
}
