# web.go

`web.go` is some kind of tiny HTML framework, templater, ... . It allows you to kind of write as little javascript as you want, and as much Go as you prefer (or any language that supports grpc I guess). Here's an example:

I have the following html template, it contains a few variables that are defined between the curly brackets. 
There's also a custom list defined before the div (which is the start of your actual template), this is a place to write pre-defined data structures. 
These will be dynamicly filled in throughout the request, this comes with the default option of retrieving your variables concurrently. 
So depending on your usecase you can have as many as you'd like. Since those define the functions that are executed (shown under html snippet).

```html
-user
 <div>
     <p>This webpage has been visited this amount of times: {counter}.
     Your last visit: {user.last_visit}</p>
     <hr>
     <p>Note if you are a human being:
     Do you wish your visits to be deleted?</p>

    <button type="button" onclick="remove({user.id})">Click Me!</button>	
</div>

<script>
    {remove}
</script>
```

```go
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
```
