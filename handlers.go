package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	utilErrors "github.com/clintjedwards/goto/errors"
	"github.com/clintjedwards/goto/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (app *app) listLinksHandler(w http.ResponseWriter, _ *http.Request) {
	links, err := app.storage.GetAllLinks()
	if err != nil {
		log.Error().Err(err).Msg("error retrieving links")
		sendErrResponse(w, http.StatusBadGateway, err)
		return
	}

	sendResponse(w, http.StatusOK, links)
}

func (app *app) createLinkHandler(w http.ResponseWriter, req *http.Request) {
	proposedLink := models.CreateLinkRequest{}

	err := parseJSON(req.Body, &proposedLink)
	if err != nil {
		log.Warn().Err(err).Msg("could not parse json")
		sendErrResponse(w, http.StatusBadRequest, err)
		return
	}
	req.Body.Close()

	err = proposedLink.Validate(app.config.MaxIDLength, req.Host)
	if err != nil {
		log.Error().Err(err).Msg("id or url invalid")
		sendErrResponse(w, http.StatusBadRequest, err)
		return
	}

	newLink := proposedLink.ToLink()

	err = app.storage.CreateLink(newLink)
	if err != nil {
		if errors.Is(err, utilErrors.ErrExists) {
			sendErrResponse(w, http.StatusConflict, err)
			return
		}
		sendErrResponse(w, http.StatusNotFound, err)
		return
	}

	log.Info().Interface("link", newLink).Msg("created new link")
	sendResponse(w, http.StatusCreated, newLink)
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

	link, err := app.storage.GetLink(linkID)
	if err != nil {
		if errors.Is(err, utilErrors.ErrNotFound) {
			sendErrResponse(w, http.StatusNotFound, err)
			return
		}
		log.Error().Err(err).Msg("error retrieving link")
		sendErrResponse(w, http.StatusBadGateway, err)
		return
	}

	returnedLink := ""

	if link.Kind == models.Formatted {
		returnedLink = generateFormattedLink(req.RequestURI[1:], link.URL)
	} else {
		returnedLink = link.URL + req.RequestURI[len(linkID)+1:]
	}

	// We wrap this so we can spit out the error to logs
	go func() {
		err := app.storage.BumpHitCount(link.ID)
		if err != nil {
			log.Error().Err(err).Msg("could not increment hit count")
		}
	}()

	http.Redirect(w, req, returnedLink, http.StatusMovedPermanently)
}

func generateFormattedLink(input string, link string) string {
	splitInput := strings.Split(input, "/")
	for _, section := range splitInput[1:] {
		link = strings.Replace(link, "{}", section, 1)
	}
	return link
}

func (app *app) getLinkHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	link, err := app.storage.GetLink(vars["id"])
	if err != nil {
		if errors.Is(err, utilErrors.ErrNotFound) {
			sendErrResponse(w, http.StatusNotFound, err)
			return
		}
		log.Error().Err(err).Msg("error retrieving link")
		sendErrResponse(w, http.StatusBadGateway, err)
		return
	}

	sendResponse(w, http.StatusOK, link)
}

func (app *app) deleteLinksHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	err := app.storage.DeleteLink(vars["id"])
	if err != nil {
		log.Error().Err(err).Msg("could not delete link")
		sendErrResponse(w, http.StatusBadGateway, err)
		return
	}

	log.Info().Str("id", vars["id"]).Msg("deleted link")
	sendResponse(w, http.StatusOK, nil)
}

// sendResponse converts raw objects and parameters to a json response
// and passes it to a provided writer.
func sendResponse(w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(payload)
	if err != nil {
		log.Error().Err(err).Msgf("could not encode json response: %v", err)
	}
}

// sendErrResponse converts raw objects and parameters to a json response specifically for erorrs
// and passes it to a provided writer. The creation of a separate function for just errors,
// is due to how they are handled differently from other payload types.
func sendErrResponse(w http.ResponseWriter, httpStatusCode int, appErr error) {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(map[string]string{"err": appErr.Error()})
	if err != nil {
		log.Error().Err(err).Msgf("could not encode json response: %v", err)
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
