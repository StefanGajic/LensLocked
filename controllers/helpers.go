package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, destination interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, destination)
}

func parseURLParams(r *http.Request, destination interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, destination)
}

func parseValues(values url.Values, destination interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(destination, values); err != nil {
		return err
	}
	return nil
}

// func parseValues(values url.Values, dst interface{}) error {
// 	dec := schema.NewDecoder()
// 	dec.IgnoreUnknownKeys(true)
// 	if err := dec.Decode(dst, values); err != nil {
// 		return err
// 	}
// 	return nil
