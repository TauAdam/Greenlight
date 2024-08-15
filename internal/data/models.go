package data

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound custom error for case when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type MockModels struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
	}
}

//func NewMockModels() MockModels {
//	return MockModels{
//		Movies: MockMovieModel{},
//	}
//}

type Models struct {
	Movies MovieModel
	Users  UserModel
}

// NewModels returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{db},
	}
}
