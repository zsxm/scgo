package chttp

import (
	"net/http"

	"github.com/zsxm/scgo/websocket"
)

type WebSocketRoute struct {
	handler func(*websocket.Conn)
}

//WebSocket路由实现ServeHTTP
func (this *WebSocketRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hand := websocket.Handler(this.handler)
	hand.ServeHTTP(w, r)
}
