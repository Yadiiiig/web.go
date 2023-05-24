package main

import (
	"webmaker/library"
)

func main() {
	library.GenFile("http://localhost:8080/foo", library.GenRequest("remove", "POST", "http://localhost:8080/foo/remove", "follow", []string{"user.id", "user.last_visit"}))
}
