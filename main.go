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
	flag.Parse()

	http.HandleFunc("/mutate", MutateCronjobs)
	http.HandleFunc("/health", Healthcheck)

	log.Print("Listening on port 8443...")

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), *cert, *key, nil))
}

func MutateCronjobs(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "%s", err)
		if err != nil {
			log.Errorf("error writing response: %v", err)
		}
	}
	mutated, err := Mutate(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "%s", err)
		if err != nil {
			log.Errorf("error writing response: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(mutated)
	if err != nil {
		log.Errorf("error writing response: %v", err)
	}
}

func Healthcheck(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
