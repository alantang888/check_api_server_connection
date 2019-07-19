package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var apiServerAddr string
var dnsTestDomain string

const API_SERVER_URL = "https://%s/api/v1"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if apiServerAddr != "" {
		_, err := http.Get(fmt.Sprintf(API_SERVER_URL, apiServerAddr))
		if err != nil {
			message := fmt.Sprintf("Connect API server error. Message: %s\n", err.Error())
			log.Println(message)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, message)
			return
		}
	}

	if dnsTestDomain != "" {
		_, err := net.LookupIP(dnsTestDomain)
		if err != nil {
			message := fmt.Sprintf("Query DNS %s error. Message: %s\n", dnsTestDomain, err.Error())
			log.Println(message)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, message)
			return
		}
	}

	log.Println("ok")
	fmt.Fprintln(w, "OK")
}

func main() {
	apiServerAddr = os.Getenv("API_SERVER_ADDR")
	dnsTestDomain = os.Getenv("DNS_TEST_DOMAIN")

	if apiServerAddr == "KUBERNETES_SERVICE_HOST" {
		apiServerAddr = os.Getenv("KUBERNETES_SERVICE_HOST")
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log.Printf("Application started. API_SERVER_ADDR: %s, DNS_TEST_DOMAIN:%s\n", apiServerAddr, dnsTestDomain)

	applicationHttp := http.NewServeMux()
	applicationHttp.HandleFunc("/health", healthHandler)
	err := http.ListenAndServe(":8080", applicationHttp)
	if err != nil {
		log.Panicf("Application HTTP server error: %s", err.Error())
	}
}
