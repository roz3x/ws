# ws
minimal websocket (ws only) lib 

# how to use 


``` go

import "github.com/roz3x/ws"

func sock(w http.ResponseWriter , r *http.Request){
	read , write := ws.Ws(w,r)
	write <- []byte("hello world")
	for {
		msg := <-read

		println("msg .. ",string(msg))
	}
}

```
