package data

import (

"database/sql"
"errors"
)
// custom error for case when
// looking up a movie that doesn't exist in our database.
var (
ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
Movies MovieModel
}
// For ease of use, we also add a container  containing
// the initialized MovieModel
func NewModels(db *sql.DB) Models {
return Models{
Movies: MovieModel{DB: db},
}
}