package main

import (
	"fmt"
	"net/http"
	"time"
)

var allServers []int

func LookActiveServer(mp *map[int]bool) {
	for {
		fmt.Println("periodic server search")
		st := 3000
		fails := 0

		for fails != 2 {
			str := fmt.Sprintf("http://localhost:%d/health", st)
			_, err := http.Get(str)
			if err != nil {
				_, ok := (*mp)[st]
				if ok {
					delete((*mp), st)
				}
				fails++
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
		if len(*mp) != 0 {
			for k := range *mp {
				allServers = append(allServers, k)
			}
		}
		time.Sleep(1 * time.Second)
		if len(*mp) != 0 {
			allServers = []int{}
		}
	}
}

func GetServer(curr *int) int {
	fmt.Println(*curr)
	*curr++
	if *curr >= len(allServers) {
		*curr = 0
	}
	dick := allServers[*curr]
	return dick
}

func main() {

	var curr int = -1

	mp := make(map[int]bool)

	go LookActiveServer(&mp)
	go GetAllServer(&mp)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		robin := GetServer(&curr)
		str := fmt.Sprintf("http://localhost:%d", robin)
		proxyReq, err := http.NewRequest(r.Method, str, r.Body)
		if err != nil {
			fmt.Println("Error occured in handler ", err)
		}
		proxyReq.Header.Set("Host", r.Host)
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

		for header, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(header, value)
			}
		}

		client := &http.Client{}
		res, err := client.Do(proxyReq)
		fmt.Println(res.Body)
		w.Write([]byte("wtf"))
		defer res.Body.Close()

	})

	http.ListenAndServe(":5173", nil)
}
