package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

/**
Please note Start functions is a placeholder for you to start your own solution.
Feel free to drop gorilla.mux if you want and use any other solution available.

main function reads host/port from env just for an example, flavor it following your taste
*/

// Start /** Starts the web server listener on given host and port.
func Start(host string, port int) {
	router := mux.NewRouter()
	router.HandleFunc("/name/{PARAM}", getNameHandler).Methods("GET")
	router.HandleFunc("/bad", getBadHandler).Methods("GET")
	router.HandleFunc("/data", getDataHandler).Methods("POST")
	router.HandleFunc("/headers", getHeadersHandler).Methods("POST")
	router.HandleFunc("/", notDefinedHandler)
	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router); err != nil {
		log.Fatal(err)
	}
}

func getNameHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hello, " + vars["PARAM"] + "!"))
}

func getBadHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
}

func getDataHandler(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err == nil {
		response := "I got message:\n" + string(body)
		writer.Write([]byte(response))
	}
}

func getHeadersHandler(writer http.ResponseWriter, request *http.Request) {
	headers := request.Header

	a, okA := headers["A"]
	b, okB := headers["B"]

	if !okA || !okB {
		getBadHandler(writer, request)
		return
	}

	aInt, errA := strconv.Atoi(a[0])
	if errA != nil {
		writeError(writer, request, errA)
		return
	}

	bInt, errB := strconv.Atoi(b[0])
	if errB != nil {
		writeError(writer, request, errB)
		return
	}

	writer.Header().Set("a+b", strconv.Itoa(aInt+bInt))
	writer.WriteHeader(http.StatusOK)
}

func notDefinedHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func writeError(writer http.ResponseWriter, request *http.Request, err error) {
	getBadHandler(writer, request)
	writer.Write([]byte(err.Error()))
}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
