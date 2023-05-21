package library

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func GenerateEndpoints(file *File) {
	http.HandleFunc(fmt.Sprintf("/%s", file.Name), func(w http.ResponseWriter, r *http.Request) {
		amount := len(file.Internal.Functions)
		results := make(chan map[int]interface{}, amount)

		var wg sync.WaitGroup
		wg.Add(amount)

		for i := 0; i < amount; i++ {
			go func(it int) {
				res, err := file.Internal.Functions[it].Run()
				if err != nil {
					log.Println(err)
				}

				results <- res

				wg.Done()
			}(i)
		}

		wg.Wait()
		close(results)

		fm := [3]interface{}{5, "20 April", "<submit>button</submit>"}
		snippet := fmt.Sprintf(file.Internal.Formatted, fm[:]...)

		fmt.Fprintf(w, snippet)
	})

}
