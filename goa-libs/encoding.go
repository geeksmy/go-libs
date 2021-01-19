package goalibs

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"mime"
	"net/http"

	"github.com/ajg/form"
	goahttp "goa.design/goa/v3/http"
)

// RequestDecoder returns a HTTP request body decoder suitable for the given
// request. The decoder handles the following mime types:
//
//     * application/json using package encoding/json
//     * application/xml using package encoding/xml
//     * application/gob using package encoding/gob
//
// RequestDecoder defaults to the JSON decoder if the request "Content-Type"
// header does not match any of the supported mime type or is missing
// altogether.
func RequestDecoder(r *http.Request) goahttp.Decoder {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		// default to JSON
		contentType = "application/json"
	} else {
		// sanitize
		if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
			contentType = mediaType
		}
	}

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body)
	case "application/gob":
		return gob.NewDecoder(r.Body)
	case "application/xml":
		return xml.NewDecoder(r.Body)
	case "application/x-www-form-urlencoded":
		return form.NewDecoder(r.Body)
	default:
		return json.NewDecoder(r.Body)
	}
}
