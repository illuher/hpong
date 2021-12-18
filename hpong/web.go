package hpong

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"inet.af/netaddr"
)

type responce struct {
	ValsList map[string]string
}

var portFlag int
var ipFlagRaw string
var ipFlag netaddr.IP

var opsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "myapp_processed_ops_total",
	Help: "The total number of processed events",
}, []string{"method", "path", "statuscode"})

func init() {
	// Tie the command-line flag to the intervalFlag variable and
	// set a usage message.
	flag.IntVar(&portFlag, "p", 8081, "Port number for listener")
	flag.StringVar(&ipFlagRaw, "ip", "0.0.0.0", "IP address to listen on")
	tmp, err := netaddr.ParseIP(ipFlagRaw)
	if err != nil {
		panic(err)
	}
	ipFlag = tmp
}

func Run() {
	flag.Parse()
	log.Println("booting")
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	pr := prometheus.NewRegistry()
	pr.MustRegister(opsProcessed)

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/ping", handlerPing)
	r.HandleFunc("/headers", handler)
	r.Handle("/metrics", promhttp.HandlerFor(pr, promhttp.HandlerOpts{}))
	r.NotFoundHandler = http.HandlerFunc(errorHandler)

	log.Println(fmt.Sprint(ipFlag) + ":" + fmt.Sprint(portFlag))

	srv := &http.Server{
		Addr: fmt.Sprint(ipFlag) + ":" + fmt.Sprint(portFlag),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	log.Println("serving")
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp := new(responce)
	resp.ValsList = make(map[string]string)

	for name, values := range r.Header {
		for _, value := range values {
			resp.ValsList[name] = value
		}
	}
	resp.ValsList["Host"] = r.Host
	resp.ValsList["IP"] = r.RemoteAddr
	opsProcessed.With(prometheus.Labels{"method": r.Method, "path": r.RequestURI, "statuscode": strconv.Itoa(200)}).Inc()
	writeResponce(w, *resp)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	resp := new(responce)
	resp.ValsList = make(map[string]string)
	resp.ValsList["message"] = "pong"
	writeResponce(w, *resp)
}

func writeResponce(w http.ResponseWriter, resp responce) {
	jsonData, err := json.Marshal(resp.ValsList)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(jsonData)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	resp := new(responce)
	resp.ValsList = make(map[string]string)
	resp.ValsList["error"] = "404"
	writeResponce(w, *resp)
}
