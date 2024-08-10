package data 

import "database/sql" 

type Movie struct {
ID int64
CreatedAt time.Time
Title string
Year int32
Runtime int32
Genres []string
Version int32
}

// MarshalJSON() method satisfies the
// json.Marshaler interface.
func (m Movie) MarshalJSON() ([]byte, error) {
var runtime string

if m.Runtime != 0 {
runtime = fmt.Sprintf("%d mins", m.Runtime)
}

aux := struct {
ID int64 `json:"id"`
Title string `json:"title"`
Year int32 `json:"year,omitempty"`
Runtime string `json:"runtime,omitempty"` // This is a string.
Genres []string `json:"genres,omitempty"`
Version int32 `json:"version"`
}{

ID: m.ID,
Title: m.Title,
Year: m.Year,
Runtime: runtime,
Genres: m.Genres,
Version: m.Version,
}

// Encode the anonymous struct to JSON, 
return json.Marshal(aux)
}


type MovieModel struct {

DB *sql.DB
}
// a placeholder method
func (m MovieModel) Insert(movie *Movie) error {
return nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
return nil, nil
}

func (m MovieModel) Update(movie *Movie) error {
return nil
}

func (m MovieModel) Delete(id int64) error {
return nil
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {

// Mock the action...
}
func (m MockMovieModel) Get(id int64) (*Movie, error) {
// Mock the action...
}
func (m MockMovieModel) Update(movie *Movie) error {
// Mock the action...
}
func (m MockMovieModel) Delete(id int64) error {
// Mock the action...
}