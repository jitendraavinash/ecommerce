package router

import (
	"fmt"
	"net/http"
)

func ServerStatus(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Server is up\n")
}
