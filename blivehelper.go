package blivehelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

type BLiveHelper struct {
	conn          *WebConn
	realId        int
	webSocketAddr string
}

type WebConn struct {
	conn *websocket.Conn
	mu   sync.RWMutex
}

func Default() *BLiveHelper {
	conn := new(WebConn)
	return &BLiveHelper{conn: conn}
}

func (b *BLiveHelper) Login(id string) (err error) {
	notlive := false
	realId, err := GetRealId(id)
	if err != nil {
		if err.Error() == ErrNOTLIVING.Error() {
			notlive = true
		} else {
			return err
		}
	}
	b.realId, err = strconv.Atoi(realId)
	if err != nil {
		return
	}
	b.webSocketAddr = GetWebSockerAddr(strconv.Itoa(b.realId))
	b.conn.conn, _, err = websocket.DefaultDialer.Dial(b.webSocketAddr, nil)
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	m["uid"] = 0
	m["roomid"] = b.realId
	m["protover"] = 1
	m["platform"] = "web"
	m["clintver"] = "1.4.0"
	data, err := json.Marshal(&m)
	if err != nil {
		return
	}
	body, err := encode(7, data)
	if err != nil {
		return
	}
	if err = b.conn.conn.WriteMessage(websocket.BinaryMessage, body); err != nil {
		return
	}
	go b.heart(20)
	log.Println("登陆成功！")
	if notlive {
		return ErrNOTLIVING
	} else {
		return nil
	}
}

func (b *BLiveHelper) GetDanmu() (chan []string, error) {
	ch := make(chan []string)

	pass := true
	go func() {
		for {
			b.conn.mu.RLock()
			_, data, err := b.conn.conn.ReadMessage()
			if err != nil {
				log.Println("读取弹幕失败！")
			}

			b.conn.mu.RUnlock()
			if pass {
				pass = false
				continue
			}
			go func(data *[]byte, ch chan<- []string) {
				msg, err := Decode(*data)
				if err != nil {
					log.Println("弹幕解码失败！")
				}
				if msg[0] != '\x00' {
					return
				}
				msgs, err := MsgDecode(msg)
				if err != nil {
					log.Println("弹幕解码失败！")
				}
				var m map[string]interface{}
				for _, v := range msgs {
					j, _ := simplejson.NewJson(v)
					m, _ = j.Map()
					if m["cmd"].(string) == "DANMU_MSG" {
						ch <- []string{"1",
							fmt.Sprintf("%s: %s",
								m["info"].([]interface{})[2].([]interface{})[1].(string),
								m["info"].([]interface{})[1].(string),
							)}
					} else if m["cmd"].(string) == "SEND_GIFT" {
						ch <- []string{"2",
							fmt.Sprintf("%s: %s了价值%v的%s*%v",
								m["data"].(map[string]interface{})["uname"].(string),
								m["data"].(map[string]interface{})["action"].(string),
								m["data"].(map[string]interface{})["price"],
								m["data"].(map[string]interface{})["giftName"].(string),
								m["data"].(map[string]interface{})["num"],
							)}
					}
				}
			}(&data, ch)
		}
	}()
	return ch, nil
}
