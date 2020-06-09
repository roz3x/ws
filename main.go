
package main

import (
	ws "./ws"
	"net/http"
)


func sock(w http.ResponseWriter , r *http.Request){
	read , write := ws.Ws(w,r)
	write <- []byte("hello world")
	for {
		msg := <-read

		println("msg .. ",string(msg))
	}
}

func main() {
	http.HandleFunc("/",func(w http.ResponseWriter , r *http.Request){
		http.ServeFile(w,r,"index.html")
	})
	http.HandleFunc("/ws",sock)
	http.ListenAndServe(":8080",nil)
}
