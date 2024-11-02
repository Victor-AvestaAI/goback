package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func HandlerReadiness(w http.ResponseWriter, req *http.Request) {
	req.Header.Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandlerValidation(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	type Message struct {
		Body string `json:"body"`
	}

	type returnErr struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	msg := Message{}

	err := decoder.Decode(&msg)

	if err != nil || len(msg.Body) > 140 {
		respErr := returnErr{}
		if err != nil {
			respErr.Error = err.Error()
		} else if len(msg.Body) > 140 {
			respErr.Error = "chirp is too long"
		}

		data, err := json.Marshal(respErr)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	type returnMessage struct {
		CleanedBody string `json:"cleaned_body"`
	}

	words := strings.Split(msg.Body, " ")

	profanity := []string{"kerfuffle", "sharbert", "fornax"}

	for idx, word := range words {
		if slices.Contains(profanity, strings.ToLower(word)) {
			words[idx] = "****"
		}
	}

	cleanedMessage := strings.Join(words, " ")

	respBody := returnMessage{
		CleanedBody: cleanedMessage,
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
