package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
	web "webmaker/library"
)

func main() {
	counter = 0
	visitsMap = make(map[string]*VisitData)

	functions := []web.Function{UserFn, CounterFn}
	actions := []web.Action{RemoveAct}

	web.Start(os.Args, functions, actions)
	//	library.GenFile("http://localhost:8080/foo", library.GenRequest("remove", "POST", "http://localhost:8080/foo/remove", "follow", []string{"user.id", "user.last_visit"}))
}

type VisitData struct {
	Visits []string
}

var (
	visitsMap map[string]*VisitData
	counter   int64
	mutex     sync.Mutex
)

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
	return "counter", fmt.Sprintf("<h1>%d</h1>", atomic.AddInt64(&counter, 1)), nil
}

func RemoveAct(args map[string]interface{}, w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	return "remove", `function remove(id) {
    console.log("removing " + id);
}`, nil
}
