package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	allServers []int
	mu         sync.RWMutex
)

func LookActiveServer(mp *map[int]bool) {
	for {
		fmt.Println("periodic server search")
		st := 3000
		fails := 0

		for fails != 2 {
			str := fmt.Sprintf("http://localhost:%d/health", st)
			_, err := http.Get(str)
			if err != nil {
				if _, ok := (*mp)[st]; ok {
					delete((*mp), st)
				}
				fails++
				st++
				continue
			}

			(*mp)[st] = true
			st++
		}

		time.Sleep(2 * time.Second)
	}
}

func GetAllServer(mp *map[int]bool) {
	for {
		// all go routines can readconcurrently after placing read lock
		mu.RLock()
		if len(*mp) != 0 {
			for k := range *mp {
				allServers = append(allServers, k)
			}
			mu.RUnlock()
		}

		time.Sleep(2 * time.Second)
		allServers = []int{}
	}
}

func GetCurrentRobin(curr *int) int {
	*curr++
	fmt.Println(*curr)
	if *curr >= len(allServers) {
		*curr = 0
	}
	robin := allServers[*curr]
	return robin
}

func main() {

	var curr int = -1

	mp := make(map[int]bool)

	go LookActiveServer(&mp)
	go GetAllServer(&mp)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		robin := GetCurrentRobin(&curr)
		str := fmt.Sprintf("http://localhost:%d", robin)
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
	})

	http.ListenAndServe(":5173", nil)
}
