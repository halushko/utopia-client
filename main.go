// main.go
package main

import (
	"fmt"
	"log"
	"strconv"

	"utopia-client/page_navigator"
)

func main() {
	id, err := strconv.Atoi("16242")
	if err != nil {
		log.Fatalf("Invalid torrent ID: %v", err)
	}

	if err := page_navigator.Login(); err != nil {
		log.Fatalf("Login error: %v", err)
	}

	html, err := page_navigator.GetTorrentPage(id)
	if err != nil {
		log.Fatalf("Fetch error: %v", err)
	}
	fmt.Println(html[:800])
}
