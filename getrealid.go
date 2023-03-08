package blivehelper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/bitly/go-simplejson"
)

// 直播间不存在
var ErrNOTEXIST = errors.New("0")

// 直播间没有开播
var ErrNOTLIVING = errors.New("1")

type GetIdData struct {
	Id int `json:"room_id"`
}

type GetId struct {
	Code int       `json:"code"`
	Data GetIdData `json:"data"`
}

func GetRealId(id string) (string, error) {
	res, err := http.Get("https://api.live.bilibili.com/room/v1/Room/room_init?id=" + id)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	j, err := simplejson.NewJson(body)
	if err != nil {
		return "", err
	}
	m, err := j.Map()
	if err != nil {
		return "", err
	}
	code, _ := m["code"].(json.Number).Int64()
	if code != int64(0) {
		return "", errors.New("0")
	}
	data, _ := m["data"].(map[string]interface{})
	realId := data["room_id"].(json.Number).String()
	if data["live_status"].(json.Number).String() == "0" {
		return realId, errors.New("1")
	} else {
		return realId, nil
	}
}
