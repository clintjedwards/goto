package main

import "errors"

// ErrNotFound is returned when a reque could not be retrieved
var ErrNotFound = errors.New("entity not found")

// ErrExists is returned when an entity is already present
var ErrExists = errors.New("entity exists")
