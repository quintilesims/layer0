package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zpatrick/go-plugin-swagger/example/handlers"
	"github.com/zpatrick/go-plugin-swagger/example/movie"
)

const (
	SWAGGER_URL      = "/api/"
	SWAGGER_SPEC_URL = "/api/swagger.json"
	SWAGGER_UI_PATH  = "static/swagger-ui/dist"
)

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	dir := http.Dir(SWAGGER_UI_PATH)
	fileServer := http.FileServer(dir)
	http.StripPrefix(SWAGGER_URL, fileServer).ServeHTTP(w, r)
}

func serveSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	spec := handlers.SwaggerSpec()
	bytes, err := json.MarshalIndent(spec, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(bytes)
}

func main() {
	// serves swagger ui from 'localhost:9090/api/'
	http.HandleFunc(SWAGGER_URL, serveSwaggerUI)

	// serves apidocs.json from 'localhost:9090/api/swagger.json'
	http.HandleFunc(SWAGGER_SPEC_URL, serveSwaggerSpec)

	// add some movies to the store
	store := movie.NewMovieStore()
	store.Insert(movie.Movie{Title: "Gladiator", Year: 2000})
	store.Insert(movie.Movie{Title: "Inception", Year: 2010})
	store.Insert(movie.Movie{Title: "Pulp Fiction", Year: 1994})

	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.ListMovies(w, r, store)
			return
		case "POST":
			handlers.AddMovie(w, r, store)
			return
		default:
			http.Error(w, "", 405)
			return
		}

	})

	log.Println("Running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
