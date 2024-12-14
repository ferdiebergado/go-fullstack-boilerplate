package response

import (
	"log"
	"net/http"
)

func RenderError(w http.ResponseWriter, err *HTTPError) {
	log.Printf("Error: %v", err.Err)
	http.Error(w, err.Error(), err.Code)
}
