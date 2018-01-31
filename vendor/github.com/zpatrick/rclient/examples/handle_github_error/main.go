package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zpatrick/rclient"
)

type repository struct {
	Name string `json:"name,omitempty"`
}

type githubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (g *githubError) Error() string {
	return g.Message
}

func githubResponseReader(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 400, 404:
		var ge *githubError
		if err := json.NewDecoder(resp.Body).Decode(&ge); err != nil {
			return err
		}

		return ge
	default:
		return json.NewDecoder(resp.Body).Decode(v)
	}
}

func main() {
	client := rclient.NewRestClient("https://api.github.com", rclient.Reader(githubResponseReader))

	var repo repository
	if err := client.Get("/repos/zpatrick/invalid_repo_name", &repo); err != nil {
		text := fmt.Sprintf("Failed to get repo: %s\n", err.Error())
		if err, ok := err.(*githubError); ok {
			text += fmt.Sprintf("Checkout Github API docs at: %s\n", err.DocumentationURL)
		}

		log.Fatal(text)
	}

	log.Println(repo)
}
