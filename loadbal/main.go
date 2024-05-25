package main

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"sync"
	"time"
)

type Server struct {
	url      string
	isActive bool
}

type Servers []*Server

var (
	servers Servers = []*Server{
		&Server{url: "http://localhost:3000", isActive: true},
		&Server{url: "http://localhost:3001", isActive: true},
		&Server{url: "http://localhost:3002", isActive: true},
	}

	timeout time.Duration = time.Second * 2

	activeSvr []*Server

	mu sync.Mutex

	curr int = -1
)

func healthChecks() {
	for {
		for i, svr := range servers {
			_, err := http.Get(svr.url + "/health")
			if err != nil {
				svr.isActive = false
				if slices.Contains(activeSvr, svr) {
					fs := servers[:i]
					ss := servers[i+1:]

					activeSvr = []*Server{}
					mu.Lock()

					for _, f := range fs {
						activeSvr = append(activeSvr, f)
					}

					for _, f := range ss {
						activeSvr = append(activeSvr, f)
					}
					mu.Unlock()
				}
				continue
			}
			mu.Lock()
			svr.isActive = true
			if slices.Contains(activeSvr, svr) {
				mu.Unlock()
				continue
			}
			activeSvr = append(activeSvr, svr)
			mu.Unlock()
		}
		time.Sleep(timeout)
	}
}

func GetCurrentRobin() *Server {
	mu.Lock()
	curr++
	fmt.Println("Request forwarded to Server ", curr)
	if curr >= len(activeSvr) {
		curr = 0
	}
	robin := activeSvr[curr]
	mu.Unlock()
	return robin
}

func ForwardRequest(w http.ResponseWriter, r *http.Request, str string) {
	proxyReq, err := http.NewRequest(r.Method, str, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
	}
	proxyReq.Header = r.Header.Clone()
	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	client := &http.Client{}
	res, err := client.Do(proxyReq)

	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
	}
	defer res.Body.Close()

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)

}

func main() {

	go healthChecks()

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		robin := GetCurrentRobin()
		ForwardRequest(w, r, robin.url)
	})

	http.ListenAndServe(":5173", nil)
}
