package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/clintjedwards/go/models"
	"go.uber.org/zap"
)

func listLinks(w http.ResponseWriter, req *http.Request) {
	l := models.Link{
		OriginalURL: "test",
		ShortURL:    "hello",
	}

	err := sendJSONResponse(w, http.StatusOK, l)
	if err != nil {
		sendJSONResponse(w, http.StatusBadGateway, "houston we have a problem")
		return
	}
}

func createLink(w http.ResponseWriter, req *http.Request) {

	newLink := models.CreateLink{}

	err := parseJSON(req.Body, &newLink)
	if err != nil {
		zap.S().Warnw("could not parse json", "error", err)
		sendJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	req.Body.Close()

	err = sendJSONResponse(w, http.StatusOK, newLink)
	if err != nil {
		zap.S().Errorw("could not send JSON response", "error", err)
		sendJSONResponse(w, http.StatusBadGateway, err)
		return
	}

	zap.S().Warnw("created new link", "link", newLink)
}

func validateLink(link models.Link) error {
	rawURL, err := url.Parse(link.OriginalURL)
	if err != nil {
		zap.S().Infow("could not parse URL", "error", err, "url", link.OriginalURL)
		return err
	}

	switch rawURL.Scheme {
	case "http", "https", "mailto", "ftp":
		break
	default:
		zap.S().Infow("invalid url scheme", "error", err, "scheme", rawURL.Scheme)
		return errors.New("invalid URL scheme")
	}

	if r.Host == rawURL.Host {
		return errRedirectLoop
	}
}

// sendJSONResponse converts raw objects and parameters to a json response and passes it to a provided writer
func sendJSONResponse(w http.ResponseWriter, httpStatusCode int, payload interface{}) error {
	w.WriteHeader(httpStatusCode)

	enc := json.NewEncoder(w)
	err := enc.Encode(payload)

	return err
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

func (app *app) createMessageHandler(w http.ResponseWriter, req *http.Request) {

	newMessage := struct {
		CallbackURL string   `json:"callback_url"` //URL to send event stream of emoji usage
		ValidEmojis []string `json:"valid_emojis"` //List of emojis to alert on
		AuthToken   string   `json:"auth_token"`   //Auth token given by app to auth on callback
		Expire      int      `json:"expire"`       //Length of time messages can be tracked. Limited to 24h
	}{}

	//Validate user supplied parameters
	err = validation.Errors{
		"callback_url": validation.Validate(newMessage.CallbackURL, is.URL),
		"valid_emojis": validation.Validate(newMessage.ValidEmojis, validation.Required),
		"auth_token":   validation.Validate(newMessage.AuthToken, validation.Required),
	}.Filter()
	if err != nil {
		sendResponse(w, http.StatusBadRequest, err.Error(), true)
		return
	}

	createdMessage := app.createMessage(newMessage.CallbackURL, newMessage.AuthToken, newMessage.ValidEmojis)

	response := struct {
		MessageID string `json:"message_id"`
	}{createdMessage.ID}

	sendResponse(w, http.StatusCreated, response, false)
	return
}
