package models

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Link is a representation of a shortened URL
type Link struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	Created int64  `json:"created"` // epoch time
	Hits    int64  `json:"hits"`    // number of visits to link
}

// Validate checks URL and Name to make sure they are valid and conform to standards
func (link Link) Validate(maxlength int) error {
	err := validation.ValidateStruct(&link,
		// URL must not be empty and a valid URL
		validation.Field(&link.URL, validation.Required, is.URL),
		// Name cannot be empty, the length must be below configured max, and must be in correct format
		validation.Field(&link.Name,
			validation.Required, validation.Length(1, maxlength), validation.By(checkValidName)),
	)
	if err != nil {
		return err
	}

	return nil
}

// checkValidName checks for a valid short name
// a name can only comprise of AlphaNumeric characters and + or _
func checkValidName(value interface{}) error {

	s, _ := value.(string)
	nameRegEx := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !nameRegEx.MatchString(s) {
		return errors.New("name is restricted to alphanumeric characters, dashes, and underscores only")
	}

	return nil
}
