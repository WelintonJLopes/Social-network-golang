package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func init() {

}

func main() {
	config.Load()
	fmt.Println(config.StringConnectionDatabase)

	fmt.Printf("Escutando na port %d", config.Port)
	r := router.Generate()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
