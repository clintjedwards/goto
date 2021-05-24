package models

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Kind represents the different types of short links.
type Kind string

const (
	// Standard links are normal short links. They take anything after the first '/' and preserve them
	// for the expanded URL.
	Standard Kind = "standard"

	// Formatted links are specialized short links that perform substitution of a user's input
	// For example:
	//
	// A user would create a link(`github`) with the URL: github.com/clintjedwards/repos/{}/issues
	// This would enable a user to write go/github/test
	// and get back this URL: github.com/clintjedwards/repos/test/issues
	// This enables a user to subtitute variables within the middle of a potentially complex URL.
	Formatted Kind = "formatted"
)

// CreateLinkRequest is a representation of the user input from a newly created link.
type CreateLinkRequest struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Link is a representation of a shortened URL
type Link struct {
	ID      string `json:"id"` // the short name of a link
	URL     string `json:"url"`
	Created int64  `json:"created"` // epoch time
	Hits    int64  `json:"hits"`    // number of visits to link
	Kind    Kind   `json:"kind"`
}

func (l CreateLinkRequest) ToLink() *Link {
	kind := Standard
	if isFormattedLink(l.URL) {
		kind = Formatted
	}

	return &Link{
		ID:      l.ID,
		URL:     l.URL,
		Created: time.Now().Unix(),
		Hits:    0,
		Kind:    kind,
	}
}

func isFormattedLink(url string) bool {
	return strings.Contains(url, "{}")
}

// Validate checks URL and ID to make sure they are valid and conform to standards
func (l CreateLinkRequest) Validate(maxlength int, serverHost string) error {
	err := validation.ValidateStruct(&l,
		// URL must not be empty and a valid URL
		validation.Field(&l.URL, validation.Required, is.URL),
		// ID cannot be empty, the length must be below configured max, and must be in correct format
		validation.Field(&l.ID,
			validation.Required, validation.Length(1, maxlength), validation.By(checkValidID)),
	)
	if err != nil {
		return err
	}

	url, _ := url.Parse(l.URL)
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
