package response

import (
	"log"
	"net/http"
)

func RenderServerError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	http.Error(w, "An error occurred.", http.StatusInternalServerError)
}
