package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

func main() {
	for i := 0; i < 500; i++ {
		go func() {
			for {
				start := time.Now()
				res, err := http.Get(os.Args[1])
				if err != nil {
					log.Println(err)
					time.Sleep(time.Second * 5)
					continue
				}
				if res.StatusCode != 200 {
					log.Println(res.StatusCode, time.Since(start))
				}

				res.Body.Close()
				incr()
			}
		}()
	}
	go printRate()
	select {}
}

var (
	lastSecond int64
	lastMinute = make([]int64, 60)
	total      int64
)

func incr() {
	atomic.AddInt64(&lastSecond, 1)
	atomic.AddInt64(&total, 1)
}

func printRate() {
	for {
		time.Sleep(time.Second * 1)
		lastMinute = append(lastMinute[1:], lastSecond)
		fmt.Printf("total: %d minute: %d second: %d\n", total, sum(lastMinute), lastSecond)
		lastSecond = 0
	}
}

func sum(sl []int64) int64 {
	n := int64(0)
	for _, s := range sl {
		n += s
	}
	return n
}
