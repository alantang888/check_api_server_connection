package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var httpTestUrl string
var dnsTestDomain string
var exitOnError bool

const API_SERVER_URL = "https://%s/api/v1"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if httpTestUrl != "" {
		_, err := http.Get(httpTestUrl)
		if err != nil {
			message := fmt.Sprintf("Connect API server error. Message: %s\n", err.Error())
			log.Println(message)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, message)
			if exitOnError {
				os.Exit(1)
			}
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
			if exitOnError {
				os.Exit(1)
			}
			return
		}
	}

	log.Println("ok")
	fmt.Fprintln(w, "OK")
}

func main() {
	httpTestUrl = os.Getenv("HTTP_TEST_URL")
	dnsTestDomain = os.Getenv("DNS_TEST_DOMAIN")

	if httpTestUrl == "KUBERNETES_SERVICE_HOST" {
		httpTestUrl = fmt.Sprintf(API_SERVER_URL, os.Getenv("KUBERNETES_SERVICE_HOST"))
	}

	if os.Getenv("EXIT_ON_ERROR") == "TRUE" {
		exitOnError = true
	} else {
		exitOnError = false
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	log.Printf("Application started. API_SERVER_ADDR: %s, DNS_TEST_DOMAIN:%s\n", httpTestUrl, dnsTestDomain)

	applicationHttp := http.NewServeMux()
	applicationHttp.HandleFunc("/health", healthHandler)
	err := http.ListenAndServe(":8080", applicationHttp)
	if err != nil {
		log.Panicf("Application HTTP server error: %s", err.Error())
	}
}
