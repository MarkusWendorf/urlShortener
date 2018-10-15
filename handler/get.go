package handler

import (
	"net/http"
	"path"
	"urlShortener/database"
)

func GetURL(db database.Database) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		short := path.Base(r.RequestURI)

		redirectTo, err := db.GetURL(short)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, redirectTo, http.StatusMovedPermanently)
		}
	}

}

