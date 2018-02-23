package model

import "errors"

func ValidateEntry(e *Entry) []error {
	var errs []error

	if e.Title == "" {
		errs = append(errs, errors.New("Title is empty"))
	}
	if e.URL == "" {
		errs = append(errs, errors.New("URL is empty"))
	}

	return errs
}
