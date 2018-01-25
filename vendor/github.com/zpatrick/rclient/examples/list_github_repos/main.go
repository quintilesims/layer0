package main

import (
	"flag"
	"fmt"
	"github.com/zpatrick/rclient"
	"log"
)

type repository struct {
	Name string `json:"name"`
}

func main() {
	username := flag.String("u", "zpatrick", "username for your github account")
	flag.Parse()

	client := rclient.NewRestClient("https://api.github.com")

	var repos []repository
	path := fmt.Sprintf("/users/%s/repos", *username)
	if err := client.Get(path, &repos); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Repos for %s: \n", *username)
	for _, r := range repos {
		fmt.Println(r.Name)
	}
}
