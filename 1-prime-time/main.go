package main

import (
	"encoding/json"
	"net"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	var debug bool
	args := os.Args
	if len(args) > 1 {
		debug = true
		log.SetLevel(log.DebugLevel)
	}
	log.SetReportTimestamp(true)
	if debug {
		log.SetReportCaller(true)
	}
	ln, err := net.Listen("tcp4", ":10000")
	if err != nil {
		log.Fatal(err)
	}
	addr := ln.Addr().String()
	log.Info("Listening on ","addr" ,addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

type data struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func handle(conn net.Conn) {
	defer conn.Close()
	var data data
	radr := conn.RemoteAddr().String()
	log.Info("Connected from %s", radr)

	err := json.NewDecoder(conn).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode the json", "decoder error", err)
	}
	log.Info("data ->","method", data.Method, "number", data.Number)
}
