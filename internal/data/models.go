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
Movies interface {
Insert(movie *Movie) error
Get(id int64) (*Movie, error)
Update(movie *Movie) error
Delete(id int64) error
}
}

func NewMockModels() Models {
return Models{
Movies: MockMovieModel{},
}
}