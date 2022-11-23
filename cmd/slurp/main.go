package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"path"

	"github.com/google/uuid"
	"moul.io/http2curl"
)

var outDir string

var upstreams arrayFlags

func main() {
	flag.StringVar(&outDir, "output-dir", "/var/lib/slurp", "directory to output the dumps to")
	flag.Var(&upstreams, "upstream", "Specify multiple times for multiple upstreams")

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

	var b bytes.Buffer
	b.ReadFrom(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	for _, upstream := range upstreams {
		err = ProxyRequest(r, upstream)
		r.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
		if err != nil {
			fmt.Printf("ERROR: can't proxy to %s\n", err.Error())
		}
	}

	w.Write([]byte("OK"))
}

func ProxyRequest(req *http.Request, destination string) error {
	fmt.Printf("proxying %s to %s\n", req.Header.Get("X-Request-ID"), destination)

	req.Host = destination
	req.URL.Host = destination
	req.URL.Scheme = "https"
	req.RequestURI = ""

	_, err := http.DefaultClient.Do(req)

	return err
}
