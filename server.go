package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"
)

var (
	config Config
	logFilter       *logutils.LevelFilter
)

func fetchTokenReq(w http.ResponseWriter, r *http.Request) {
	log.Print("[DEBUG] Received fetchToken request")

	code, ok := r.URL.Query()["code"]
	if !ok {
		http.Error(w, "Couldn't read authorization code", http.StatusBadRequest)
	}

	queryParameters := "?client_id=" + config.ClientID + "&client_secret=" + config.ClientSecret + "&code=" + code[0]
	url := TokenURL() + queryParameters
	log.Printf("[DEBUG] calling %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] Request failed: %v", err)
		http.Error(w, "Couldn't fetch token", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Request failed, could not read response: %v", err)
		http.Error(w, "Couldn't get read response", resp.StatusCode)
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	logFilter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}
	log.SetOutput(logFilter)

	LoadConfig()

	router := mux.NewRouter()

	router.HandleFunc("/token", fetchTokenReq).Methods("GET")

	log.Printf("[DEBUG] listening on %s, with internal port %s","/token", config.Port)
	//log.Fatal(http.ListenAndServeTLS(":"+port, certificate, key, router))
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
