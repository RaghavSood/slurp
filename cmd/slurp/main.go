package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"path"

	"github.com/google/uuid"
	"moul.io/http2curl"
)

var outDir string

func main() {
	flag.StringVar(&outDir, "output-dir", "/var/lib/slurp", "directory to output the dumps to")

	flag.Parse()

	router := http.NewServeMux()

	router.HandleFunc("/", indexHandler)

	err := http.ListenAndServe("localhost:8000", router)
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	fmt.Println(cleanPath)
	if r.Header.Get("X-Request-ID") == "" {
		r.Header.Set("X-Request-ID", uuid.New().String())
	}

	reqId := r.Header.Get("X-Request-ID")
	dumpReq, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("DR", err)
	}

	curl, err := http2curl.GetCurlCommand(r)
	if err != nil {
		fmt.Println("curl", err)
	}

	fmt.Println("DR")
	fmt.Println(string(dumpReq))
	writeFile(outDir, cleanPath, reqId, "req", string(dumpReq))
	fmt.Println("curl")
	fmt.Println(curl)
	writeFile(outDir, cleanPath, reqId, "curl", curl.String())

	w.Write([]byte("OK"))
}
