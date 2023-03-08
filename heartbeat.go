package blivehelper

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (b *BLiveHelper) heart(heart int) {
	for {
		time.Sleep(time.Duration(heart) * time.Second)
		b.conn.mu.Lock()
		body, err := encode(2, []byte{})
		if err != nil {
			log.Println(err)
		}
		err = b.conn.conn.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			log.Println(err)
		}
		b.conn.mu.Unlock()
	}
}
