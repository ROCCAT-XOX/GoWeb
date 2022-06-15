package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Response(message string, status int, rw http.ResponseWriter) {
	rw.WriteHeader(status)
	rw.Write([]byte(message))
}

func Info(s string) {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info(s)
}
