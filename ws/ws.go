package ws

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	fin        = 1 << 7
	rsv1       = 1 << 6
	rsv2       = 1 << 5
	rsv3       = 1 << 4
	textmode   = 1
	binarymode = 2
)

var (
	buff    *bufio.ReadWriter
	newConn net.Conn
)

func Ws(w http.ResponseWriter, r *http.Request)(chan []byte , chan []byte) {
	hiJ, ok := w.(http.Hijacker)
	if !ok {
		panic("error")
	}
	conn, _, err := hiJ.Hijack()
	if err != nil {
		panic("error")
	}
	hash := sha1.New()
	io.WriteString(hash, r.Header["Sec-Websocket-Key"][0])
	io.WriteString(hash, "258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	hashsum := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	resp := make([]byte, 0)
	resp = append(resp, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	resp = append(resp, hashsum...)
	resp = append(resp, "\r\n\r\n"...)
	_, err = conn.Write(resp)
	if err != nil {
		panic("error")
	}
	conn.SetDeadline(time.Time{})



	write := make( chan []byte)
	read  := make( chan []byte)

	go func(){
		for {
			msg :=  <-write
			println(string(msg))
			m := make([]byte,0,1024)
			m = append(m,byte(fin|textmode),byte(len(msg)))
			m = append(m,msg...)
			conn.Write(m)
		}
	}()
	go func(){
		for {
			msg := make([]byte,1024)
			_ , _ = conn.Read(msg)
			read <- msg
		}
	}()
	return read , write
}

