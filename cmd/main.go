package main

import (
	"duplicate/server"
	"log"
)


func main() {
	err := server.Run()	
	if err != nil {
		log.Fatal(err)
	}
}