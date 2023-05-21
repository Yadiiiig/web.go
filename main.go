package main

import (
	"log"
	web "webmaker/library"
)

func main() {
	err := web.Scan(".files", User, Counter)
	if err != nil {
		log.Fatal(err)
	}
}

func User() (map[int]interface{}, error) {
	return nil, nil
}

func Counter() (map[int]interface{}, error) {
	return nil, nil
}
