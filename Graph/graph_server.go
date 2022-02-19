package Graph

import (
	"context"
	"net/http"
	"time"
)

var server = http.Server{Addr: ":8081"}

func lineserver(w http.ResponseWriter, _ *http.Request) {
	CreateLineGraph(w)
}
func pieserver(w http.ResponseWriter, _ *http.Request) {
	CreatePieGraph(w)
}

func StartServer() {
	mux := http.NewServeMux()
	server.Handler = mux
	mux.HandleFunc("/line", lineserver)
	mux.HandleFunc("/pie", pieserver)
	server.ListenAndServe()
}

func StopSever(){
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	server.Shutdown(ctx)
	//Close other stuff, if needed
	cancel()
}