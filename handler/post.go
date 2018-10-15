package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"urlShortener/database"
)

type Request struct {
	Url string `json:"url"`
}

func PostShorthand(db database.Database) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		var request Request
		defer r.Body.Close()

		body := bytes.Buffer{}
		io.Copy(&body, r.Body)

		err := json.Unmarshal(body.Bytes(), &request)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, err)
			return
		}

		var shorthand string
		for {
			// generate new shorthands until unique (not in database yet)
			shorthand = generateShorthand()
			_, err := db.GetURL(shorthand)
			if err == database.ErrShorthandDoesNotExist {
				break
			}
		}

		response, err := json.Marshal(Request{Url: shorthand})
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		err = db.SetURL(shorthand, request.Url)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, err)
			return
		}

		jsonResponse(w, http.StatusOK, response)
	}
}


func generateShorthand() string {

	idLen := 8
	alphabet := "abcdefghijklmopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXYZ"
	alphabetLen := byte(len(alphabet))

	random := make([]byte, idLen)
	rand.Reader.Read(random)

	id := make([]byte, idLen)
	for i := 0; i < idLen; i++ {
		id[i] = alphabet[random[i]%alphabetLen]
	}

	return string(id)
}

func jsonResponse(w http.ResponseWriter, status int, json []byte) {

	w.WriteHeader(status)
	headers := w.Header()
	headers.Add("Content-Type", "application/json")

	w.Write(json)
}

func errorResponse(w http.ResponseWriter, status int, err error) {

	if err == nil {
		log.Fatal("parameter err of type error is nil")
	}

	jsonResponse(w, status, []byte("{err:" + err.Error() + "}"))
}