package main

import (
	"encoding/json"
	"io"
	"math"
	"net"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportTimestamp(true)
	log.SetReportCaller(true)
	ln, err := net.Listen("tcp4", ":10000")
	if err != nil {
		log.Fatal(err)
	}
	addr := ln.Addr().String()
	log.Info("Listening on ", "addr", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}

type request struct {
	Method string      `json:"method"`
	Number interface{} `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func handle(conn net.Conn) {
	radr := conn.RemoteAddr().String()
	log.Info("New connection ->", "Connected from", radr)
	for {
		var data request
		var r response
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(&data)
		if err != nil {
			if err != io.EOF {
				if _, ok := err.(*json.SyntaxError); ok {
					log.Warn("malformed request , failed to decode json", "decoder error", err, "")
					buf, _ := io.ReadAll(decoder.Buffered())
					log.Warn("raw request", "buf", string(buf))
				}
			}
		}
		log.Info("data ->", "method", data.Method, "number", data.Number)

		if data.isValid() { // request is well formed
			log.Info("request is VALID", "req", data)
			if data.hasPrime() { // number is pime
				log.Info("number IS PRIME", "num", data.Number)
				r.Method = "isPrime"
				r.Prime = true
				log.Info("sending response ", "resp", r)
				if err := json.NewEncoder(conn).Encode(r); err != nil {
					log.Fatal("failed to encode response", "err", err)
				}
			} else { // not prime
				log.Info("number is NOT PRIME", "num", data.Number)
				r.Method = "isPrime"
				r.Prime = false
				log.Info("sending response ", "resp", r)
				if err := json.NewEncoder(conn).Encode(r); err != nil {
					log.Fatal("failed to encode response", "err", err)
				}
			}
		} else { // request is malformed
			log.Info("request is MALFORMED", "req", data)
			sendMalformedRespBack(conn, data)
			break
		}
	}
	log.Warn("closing connection")
	defer conn.Close()
}

func sendMalformedRespBack(conn net.Conn, req request) {
	var r response
	r.Method = "MALFORMED"
	r.Prime = false
	log.Warn("sending response ->", "resp", r, "orignalData", req)
	if err := json.NewEncoder(conn).Encode(r); err != nil {
		log.Fatal("failed to encode response", "err", err)
	}
}

func (d *request) hasPrime() bool {
	switch num := d.Number.(type) {
	case int:
		log.Info("num is an int", "num", num)
		return isPrime(num)
	case float64:
		floor := math.Floor(num)
		if num == floor {
			log.Debug("these are equal", "num", num, "floor", floor)
			log.Info("checking float64 converted int", "orignal", num, "converted", int(num))
			return isPrime(int(num))
		}
		log.Debug("num is a float", "float", num)
		return false
	default:
		log.Warnf("%s with type %T is not a prime number", num, num)
	}
	return false
}

func isPrime(n int) bool {
	log.Info("checking if n is prime", "n", n)
	if n <= 1 {
		log.Info("Not prime", "n", n)
		return false
	}
	if n <= 3 {
		log.Info("is prime", "n", n)
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		log.Info("Not prime", "n", n)
		return false
	}
	for i := 5; i*i < n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			log.Info("Not prime", "n", n)
			return false
		}
	}
	log.Info("is prime", "n", n)
	return true
}

func (d *request) isValid() bool {
	if d.Method != "isPrime" {
		return false
	}
	switch d.Number.(type) {
	case int:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}
