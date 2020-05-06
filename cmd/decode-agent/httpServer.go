package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/yh742/j2735-decoder/pkg/decoder"
)

// ExposedSettings are settings exposed through http
type ExposedSettings struct {
	PubTopic string
	SubTopic string
	Format   decoder.StringFormatType
}

// HTTPServer is use for exposing bridge settings
type HTTPServer struct {
	router     *mux.Router
	httpHupSig chan bool
	server     http.Server
	port       int
	auth       basicAuth
}

// NewHTTPServer creates a new http server instance
func NewHTTPServer(port int) *HTTPServer {
	server := HTTPServer{
		router:     mux.NewRouter(),
		httpHupSig: make(chan bool),
		port:       port,
	}
	return &server
}

// RegisterBridge registers a new bridge connection setting for exposure
func (hs *HTTPServer) RegisterBridge(bd *bridge) {
	url := "/" + strings.ToLower(bd.cfg.Name) + "/settings"
	hs.auth = parseAuthFiles(bd.cfg.Op.HTTPAuth)
	log.Debug().Msgf("Username: '%s' Password: '%s'", hs.auth.username, hs.auth.password)

	// GET methods
	hs.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		getSettingHandler(w, r, bd.cfg, hs.auth)
	}).Methods("GET")

	// PUT methods
	hs.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		putSettingsHandler(w, r, hs.auth, func(eSetting ExposedSettings) {
			bd.updateSettings(eSetting)
		})
	}).Methods("PUT")
}

// StartListening starts the listening for requests
func (hs *HTTPServer) StartListening(block bool) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Origin", "Accept", "Authorization", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	allowCreds := handlers.AllowCredentials()
	hs.server = http.Server{
		Addr:    ":" + strconv.Itoa(hs.port),
		Handler: handlers.CORS(originsOk, headersOk, methodsOk, allowCreds)(hs.router),
	}
	runServer := func() {
		if err := hs.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed at listenandserver()")
		}
		hs.httpHupSig <- true
	}
	if (!block) {
		go runServer()
	} else {
		runServer()
	}
}

// Disconnect tearsdown the HTTP server instance
func (hs *HTTPServer) Disconnect() {
	if hs.server.Addr != "" {
		if err := hs.server.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("can't shutdown http server properly")
		}
		log.Debug().Msg("http server teardown...")
		<-hs.httpHupSig
		log.Debug().Msg("http server teardown finished...")
	}
}
