package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func main() {
	cert := flag.String("cert", "/etc/admission-webhook/tls/tls.crt", "Path to the certificate file")
	key := flag.String("key", "/etc/admission-webhook/tls/tls.key", "Path to the key file")
	port := flag.Int("port", 8443, "Port to listen on")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	http.HandleFunc("/mutate-cronjob", Mutate)
	http.HandleFunc("/mutate-job", Mutate)
	http.HandleFunc("/health", Healthcheck)

	log.Printf("Listening on port %d...", *port)

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *cert, *key, nil))
}

func Mutate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Error closing request body: %v", err)
		}
	}(r.Body)
	const writingErrorFormat = "Error writing response: %v"
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "%s", err)
		if err != nil {
			log.Errorf(writingErrorFormat, err)
		}
	}
	mutated := []byte{}
	log.Debugf("Received request: %s", r.URL.Path)
	log.Debugf("recv: %s\n", string(body))
	switch r.URL.Path {
	case "/mutate-cronjob":
		mutated, err = MutateCronjobs(body)
	case "/mutate-job":
		mutated, err = MutateJobs(body)
	}

	if err != nil {
		log.Errorf("Error mutating: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "Error mutating: %s", err)
		if err != nil {
			log.Errorf(writingErrorFormat, err)
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(mutated)
	if err != nil {
		log.Errorf(writingErrorFormat, err)
	}
}

func Healthcheck(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
