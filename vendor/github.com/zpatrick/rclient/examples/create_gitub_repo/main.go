package main

import (
	"flag"
	"fmt"
	"github.com/zpatrick/rclient"
	"log"
)

type repository struct {
	Name string `json:"name,omitempty"`
}

func main() {
	username := flag.String("u", "", "username for your github account")
	password := flag.String("p", "", "password for your github account")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("username and password are required")
	}

	client := rclient.NewRestClient("https://api.github.com")

	var repo repository
	request := repository{Name: "my_sample_repo"}

	// add auth to the request
	if err := client.Post("/user/repos", request, &repo, rclient.BasicAuth(*username, *password)); err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	fmt.Printf("Successfully created repository %s\n", repo.Name)

	// also, you can set basic auth for each request the client makes
	client = rclient.NewRestClient("https://api.github.com", rclient.RequestOptions(rclient.BasicAuth(*username, *password)))

	path := fmt.Sprintf("/repos/%s/%s", *username, repo.Name)
	if err := client.Delete(path, nil, nil); err != nil {
		log.Fatalf("Failed to delete repository: %v", err)
	}

	fmt.Printf("Successfully deleted repository %s\n", repo.Name)
}
