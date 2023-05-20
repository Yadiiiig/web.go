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

func User() (string, error) {
	return "", nil
}

func Counter() (string, error) {
	return "", nil
}
