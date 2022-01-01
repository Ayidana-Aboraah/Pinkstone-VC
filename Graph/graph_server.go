package Graph

import (
	"net/http"
)

func lineserver(w http.ResponseWriter, _ *http.Request) {
	CreateLineGraph(w)
}
func pieserver(w http.ResponseWriter, _ *http.Request) {
	CreatePieGraph(w)
}

func StartServers() {
	http.HandleFunc("/line", lineserver)
	http.HandleFunc("/pie", pieserver)

	http.ListenAndServe(":8081", nil)
}
