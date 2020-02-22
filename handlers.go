package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/clintjedwards/go/models"
	"go.uber.org/zap"
)

func (app *app) listLinksHandler(w http.ResponseWriter, req *http.Request) {
	links, err := app.storage.GetAllLinks()
	if err != nil {
		zap.S().Errorw("error in retrieving links", "error", err)
		sendJSONResponse(w, http.StatusBadGateway, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, links)
}

func (app *app) createLinkHandler(w http.ResponseWriter, req *http.Request) {

	newLink := models.Link{}

	err := parseJSON(req.Body, &newLink)
	if err != nil {
		zap.S().Warnw("could not parse json", "error", err)
		sendJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	req.Body.Close()

	err = newLink.Validate(app.config.MaxNameLength, req.Host)
	if err != nil {
		zap.S().Errorw("name or url invalid", "error", err)
		sendJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	newLink.Created = time.Now().Unix()

	app.storage.CreateLink(newLink)

	zap.S().Infow("created new link", "link", newLink)
	sendJSONResponse(w, http.StatusCreated, newLink)
}

// sendJSONResponse converts raw objects and parameters to a json response and passes it to a provided writer
func sendJSONResponse(w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(payload)
	if err != nil {
		zap.S().Errorw("could not send JSON response", "error", err)
	}
}

// parseJSON parses the given json request into interface
func parseJSON(rc io.Reader, object interface{}) error {
	decoder := json.NewDecoder(rc)
	err := decoder.Decode(object)
	if err != nil {
		log.Println(err)
		return errors.New("could not parse json")
	}
	return nil
}
