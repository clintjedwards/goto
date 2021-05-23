package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/clintjedwards/goto/models"
	"github.com/clintjedwards/toolkit/tkerrors"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (app *app) listLinksHandler(w http.ResponseWriter, req *http.Request) {
	links, err := app.storage.GetAllLinks()
	if err != nil {
		log.Error().Err(err).Msg("error retrieving link ")
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, links)
}

func (app *app) createLinkHandler(w http.ResponseWriter, req *http.Request) {
	proposedLink := models.CreateLinkRequest{}

	err := parseJSON(req.Body, &proposedLink)
	if err != nil {
		log.Warn().Err(err).Msg("could not parse json")
		sendJSONErrResponse(w, http.StatusBadRequest, err)
		return
	}
	req.Body.Close()

	err = proposedLink.Validate(app.config.MaxIDLength, req.Host)
	if err != nil {
		log.Error().Err(err).Msg("id or url invalid")
		sendJSONErrResponse(w, http.StatusBadRequest, err)
		return
	}

	newLink := proposedLink.ToLink()

	err = app.storage.CreateLink(newLink)
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityExists) {
			sendJSONErrResponse(w, http.StatusConflict, err)
			return
		}
		sendJSONErrResponse(w, http.StatusNotFound, err)
		return
	}

	log.Info().Interface("link", newLink).Msg("created new link")
	sendJSONResponse(w, http.StatusCreated, newLink)
}

func isReservedCharacter(c rune) bool {
	reservedChars := map[rune]struct{}{
		'!': {}, '#': {}, '$': {}, '&': {}, '\'': {}, '(': {},
		')': {}, '*': {}, '+': {}, ',': {}, '/': {}, ':': {}, ';': {}, '=': {},
		'?': {}, '@': {}, '[': {}, ']': {},
	}

	if _, ok := reservedChars[c]; ok {
		return true
	}

	return false
}

func (app *app) followLinkHandler(w http.ResponseWriter, req *http.Request) {

	splitURL := strings.FieldsFunc(req.RequestURI[1:], isReservedCharacter)
	linkID := splitURL[0]
	queryParams := req.RequestURI[len(linkID)+1:]

	link, err := app.storage.GetLink(linkID)
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityNotFound) {
			sendJSONErrResponse(w, http.StatusNotFound, err)
			return
		}
		log.Error().Err(err).Msg("error retrieving link")
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	// We wrap this so we can spit out the error to logs
	go func() {
		err := app.storage.BumpHitCount(link.ID)
		if err != nil {
			log.Error().Err(err).Msg("could not increment hit count")
		}
	}()

	http.Redirect(w, req, link.URL+queryParams, http.StatusMovedPermanently)
}

func (app *app) getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	link, err := app.storage.GetLink(vars["id"])
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityNotFound) {
			sendJSONErrResponse(w, http.StatusNotFound, err)
			return
		}
		log.Error().Err(err).Msg("error retrieving link")
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, link)
}

func (app *app) deleteLinksHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	err := app.storage.DeleteLink(vars["id"])
	if err != nil {
		log.Error().Err(err).Msg("could not delete link")
		sendJSONErrResponse(w, http.StatusBadGateway, err)
		return
	}

	log.Info().Str("id", vars["id"]).Msg("deleted link")
	sendJSONResponse(w, http.StatusOK, nil)
}

// sendJSONResponse converts raw objects and parameters to a json response and passes it to a provided writer
func sendJSONResponse(w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(payload)
	if err != nil {
		log.Error().Err(err).Msg("could not send JSON response")
	}
}

// sendJSONErrResponse converts raw objects and parameters to a json response and passes it to a provided writer
func sendJSONErrResponse(w http.ResponseWriter, httpStatusCode int, errStr error) {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]string{"err": errStr.Error()})
	if err != nil {
		log.Error().Err(err).Msg("could not send JSON response")
	}
}

// parseJSON parses the given json request into interface
func parseJSON(rc io.Reader, object interface{}) error {
	decoder := json.NewDecoder(rc)
	err := decoder.Decode(object)
	if err != nil {
		return errors.New("could not parse json")
	}
	return nil
}
