package models

import (
	"errors"
	"net/url"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Link is a representation of a shortened URL
type Link struct {
	ID      string `json:"id"` // the short name of a link
	URL     string `json:"url"`
	Created int64  `json:"created"` // epoch time
	Hits    int64  `json:"hits"`    // number of visits to link
}

// Validate checks URL and ID to make sure they are valid and conform to standards
func (link Link) Validate(maxlength int, serverHost string) error {
	err := validation.ValidateStruct(&link,
		// URL must not be empty and a valid URL
		validation.Field(&link.URL, validation.Required, is.URL),
		// ID cannot be empty, the length must be below configured max, and must be in correct format
		validation.Field(&link.ID,
			validation.Required, validation.Length(1, maxlength), validation.By(checkValidID)),
	)
	if err != nil {
		return err
	}

	url, _ := url.Parse(link.URL)
	if serverHost == url.Host {
		return errors.New("redirect loop detected")
	}

	return nil
}

// checkValidID checks for a valid short name
// an ID can only comprise of AlphaNumeric characters and + or _
func checkValidID(value interface{}) error {

	reservedIDs := []string{"links", "create", "version", "status", "health", "edit", "api"}

	s, _ := value.(string)
	idRegEx := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !idRegEx.MatchString(s) {
		return errors.New("id is restricted to alphanumeric characters, dashes, and underscores only")
	}

	for _, id := range reservedIDs {
		if s == id {
			return errors.New("requested id is reserved and cannot be used")
		}
	}

	return nil
}
