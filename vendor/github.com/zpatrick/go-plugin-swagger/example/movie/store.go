package movie

type MovieStore struct {
	movies []Movie
}

func NewMovieStore() *MovieStore {
	return &MovieStore{
		movies: []Movie{},
	}
}

func (m *MovieStore) Movies() []Movie {
	return m.movies
}

func (m *MovieStore) Insert(movie Movie) {
	m.movies = append(m.movies, movie)
}
