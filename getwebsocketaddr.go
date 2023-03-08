package blivehelper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetAddrData struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type GetAddr struct {
	Code int         `json:"code"`
	Data GetAddrData `json:"data"`
}

func GetWebSockerAddr(id string) string {
	checkErr := func(err error) {
		if err != nil {
			log.Println(err)
		}
	}
	addr := []string{
		"https://api.live.bilibili.com/room/v1/Danmu/getConf?room_id=",
		"&platform=pc&player=web",
	}
	res, err := http.Get(addr[0] + id + addr[1])
	checkErr(err)
	body, err := io.ReadAll(res.Body)
	checkErr(err)
	defer res.Body.Close()

	var project GetAddr
	err = json.Unmarshal(body, &project)
	checkErr(err)
	return "wss://" + project.Data.Host + "/sub"
}
