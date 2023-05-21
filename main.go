package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	web "webmaker/library"
)

type VisitData struct {
	Visits []string
}

var (
	visitsMap map[string]*VisitData
	counter   int64
	mutex     sync.Mutex
)

func main() {
	counter = 0
	visitsMap = make(map[string]*VisitData)

	err := web.Scan(".files", UserFn, CounterFn, RemoveFn)
	if err != nil {
		log.Fatal(err)
	}
}

type UserDb struct {
	Id        int    `web:"id"`
	LastVisit string `web:"last_visit"`
}

func UserFn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	mutex.Lock()

	visitData, exists := visitsMap[ip]
	if !exists {
		visitData = &VisitData{}
		visitsMap[ip] = visitData
	}

	visitData.Visits = append(visitData.Visits, time.Now().Format(time.RFC3339))

	mutex.Unlock()

	return "user", UserDb{1, time.Now().Format(time.RFC3339)}, nil
}

func CounterFn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "counter", atomic.AddInt64(&counter, 1), nil
}

func RemoveFn(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "remove", `function remove(id) {
    console.log("removing " + id);
}`, nil
}
