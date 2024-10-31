package services

import (
	"database/sql"

	"github.com/icpinto/dating-app/repositories"
)

func GetUsepwd(username string, db *sql.DB) (string, int, error) {
	return repositories.GetUserpwdByUsername(db, username)
}
