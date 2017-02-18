package server

import (
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	HostEnvironment = "GFS_CURRENCY_HOST" // hostname environment variable
	PortEnvironment = "GFS_CURRENCY_PORT" // port environment variable

	defaultHost = "127.0.0.1" // default hostname
	defaultPort = "4000"      // default port number
)

type Server struct {
	host string
	port int

	hasCurrencies  bool               // true if currencies have been properly fetched+parsed
	lastUpdateTime time.Time          // time parsed from timestamp in ECB data
	currencies     map[string]float64 // currency data

	mutex    *sync.Mutex        // used for locking when handling webhooks
	webhooks map[string]webhook // holds webhooks

	currencyHits    *expvar.Int
	convertHits     *expvar.Int
	webhookHits     *expvar.Int
	webhookTriggers *expvar.Int
}

// representation of a webhook
type webhook struct {
	BaseCurrency string `json:"base_currency"`
	Url          string `json:"url"`
	Secret       string `json:"secret"`
}

// Creates a new server. Creation reads environment variables to configure
// hostname and port.
func New() (s *Server, err error) {
	// get config from environment
	host := getEnv(HostEnvironment, defaultHost)
	portStr := getEnv(PortEnvironment, defaultPort)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("Error parsing port number: %s", portStr)
	}

	// initialize internal variables
	return &Server{
		host: host,
		port: port,

		hasCurrencies: false,

		mutex:    &sync.Mutex{},
		webhooks: make(map[string]webhook),

		currencyHits:    expvar.NewInt("currency_hits"),
		convertHits:     expvar.NewInt("convert_hits"),
		webhookHits:     expvar.NewInt("webhook_hits"),
		webhookTriggers: expvar.NewInt("webhook_triggers"),
	}, nil
}

// Runs the server and returns error from http.ListenAndServe
func (s *Server) Run() (err error) {
	log.Printf("Starting server on %s:%d\n", s.host, s.port)

	// starts the currency updating goroutine
	s.startCurrencyUpdating()
	http.Handle("/", s)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
}

// gets the value of the given environment variable, if no value is set
// the given default is returned
func getEnv(env, def string) (val string) {
	val = os.Getenv(env)
	if val == "" {
		val = def
	}

	return val
}
