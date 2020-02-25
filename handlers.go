package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/clintjedwards/go/models"
	"github.com/clintjedwards/toolkit/tkerrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (app *app) listLinksHandler(w http.ResponseWriter, req *http.Request) {
	links, err := app.storage.GetAllLinks()
	if err != nil {
		zap.S().Errorw("error in retrieving links", "error", err)
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, links)
}

func (app *app) createLinkHandler(w http.ResponseWriter, req *http.Request) {

	newLink := models.Link{}

	err := parseJSON(req.Body, &newLink)
	if err != nil {
		zap.S().Warnw("could not parse json", "error", err)
		sendJSONErrResponse(w, http.StatusBadRequest, err)
		return
	}
	req.Body.Close()

	err = newLink.Validate(app.config.MaxNameLength, req.Host)
	if err != nil {
		zap.S().Errorw("name or url invalid", "error", err)
		sendJSONErrResponse(w, http.StatusBadRequest, err)
		return
	}

	newLink.Created = time.Now().Unix()
	newLink.Hits = 0

	err = app.storage.CreateLink(newLink)
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityExists) {
			sendJSONErrResponse(w, http.StatusConflict, err)
			return
		}
		sendJSONErrResponse(w, http.StatusNotFound, err)
		return
	}

	zap.S().Infow("created new link", "link", newLink)
	sendJSONResponse(w, http.StatusCreated, newLink)
}

func (app *app) followLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	link, err := app.storage.GetLink(vars["name"])
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityNotFound) {
			sendJSONErrResponse(w, http.StatusNotFound, err)
			return
		}
		zap.S().Errorw("error retrieving link", "error", err)
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	// We wrap this so we can spit out the error to logs
	go func() {
		err := app.storage.BumpHitCount(link.Name)
		if err != nil {
			zap.S().Errorw("could not increment hit count", "error", err)
		}
	}()

	http.Redirect(w, req, link.URL, http.StatusMovedPermanently)
}

func (app *app) getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	link, err := app.storage.GetLink(vars["name"])
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityNotFound) {
			sendJSONErrResponse(w, http.StatusNotFound, err)
			return
		}
		zap.S().Errorw("error retrieving link", "error", err)
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, link)
}

func (app *app) deleteLinksHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	err := app.storage.DeleteLink(vars["name"])
	if err != nil {
		zap.S().Errorw("could not delete link", "error", err)
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	zap.S().Infow("deleted link", "name", vars["name"])
	sendJSONResponse(w, http.StatusOK, nil)
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

// sendJSONErrResponse converts raw objects and parameters to a json response and passes it to a provided writer
func sendJSONErrResponse(w http.ResponseWriter, httpStatusCode int, errStr error) {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]string{"err": errStr.Error()})
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
