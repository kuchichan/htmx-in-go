package main

import (
	"errors"
	"net/http"

	"github.com/go-playground/form/v4"
)

func (app *application) decodeForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.decoder.Decode(dst, r.PostForm)

	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func isNotEmpty(value string) bool {
	return value != ""
}

func isPositive(value int) bool {
	return value > 0
}
