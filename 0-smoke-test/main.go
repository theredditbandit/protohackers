package main

import (
	"context"
	"io"
	"log"
	"net"
)

func main() {
	n := 1
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	l, err := net.Listen("tcp4", ":6942")
	if err != nil {
		log.Fatal(err)
	}
	addr := l.Addr().String()
	log.Printf("Listening on %s", addr)
	defer l.Close()

	for {
		log.Printf("listening for conn %d", n)
		conn, err := l.Accept()
		log.Printf("accepted conn %d", n)
		n++
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			ctx := context.WithValue(context.Background(), "clientAddr", conn.RemoteAddr().String())
			buff := make([]byte, 4096)
			for {
				n, err := conn.Read(buff)
				if err != nil {
					if err == io.EOF {
						log.Printf("Client %s has closed the connection", ctx.Value("clientAddr"))
						break
					}
					log.Fatalf("failed to read data from client %s ,\n Err : %s", ctx.Value("clientAddr"), err)
				}
				log.Printf("(%s) Got data : %s ", ctx.Value("clientAddr"), buff[:n])
				_, err = conn.Write(buff[:n])
				if err != nil {
					log.Fatalf("(%s) failed to write data Err %s:", ctx.Value("clientAddr"), err)
				}
			}
		}()
		log.Print("==========================================================")
	}
}
