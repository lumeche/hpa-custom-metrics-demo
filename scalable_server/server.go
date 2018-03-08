package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
)

func main() {
	var ops uint32
	hostname := os.Getenv("HOSTNAME")
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&ops, 1)
		time.Sleep(10 * time.Second)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		atomic.AddUint32(&ops, ^uint32(0))
	})

	http.HandleFunc("/total", func(w http.ResponseWriter, r *http.Request) {
		total := atomic.LoadUint32(&ops)
		log.Printf("hostname:%s, total:%d", hostname, total)
		fmt.Fprint(w, hostname+","+fmt.Sprint(total))
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
