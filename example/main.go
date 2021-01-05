package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
	negroniprometheus "github.com/zbindenren/negroni-prometheus"
)

func main() {
	n := negroni.New()
	m := negroniprometheus.NewMiddleware("serviceName")
	// if you want to use other buckets than the default (300, 1200, 5000) you can run:
	// m := negroniprometheus.NewMiddleware("serviceName", 400, 1600, 700)

	n.Use(m)

	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("200 - Something bad happened!"))
		fmt.Fprintf(w, "slept %d milliseconds\n", sleep)
	})
	r.HandleFunc(`/notfound`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Something bad happened!"))
		fmt.Fprintln(w, "not found")
	})

	n.UseHandler(r)

	n.Run(":3000")
}
