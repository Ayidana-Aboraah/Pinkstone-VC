package Graph

import (
	"context"
	"net/http"
	"time"
)

var server = http.Server{Addr: ":8081"}

func StartServer() {
	mux := http.NewServeMux()
	server.Handler = mux
	mux.HandleFunc("/line", func(w http.ResponseWriter, r *http.Request) { CreateLineGraph(w) })
	mux.HandleFunc("/pie", func(w http.ResponseWriter, r *http.Request) { CreatePieGraph(w) })
	server.ListenAndServe()
}

func StopSever() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(ctx)
	cancel() //Close other stuff if needed
}
