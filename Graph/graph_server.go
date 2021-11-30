package Graph

import (
	"net/http"
)

func httpserver(w http.ResponseWriter, _ *http.Request) {
	CreateBarGraph(w)
}

func StartServer() {
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
}