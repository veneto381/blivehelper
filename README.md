# blivehelper
bilibili（看）直播小工具
目前有的功能：获取直播弹幕、获取直播间礼物信息  
TODO:
- [ ] 登陆
- [ ] 发送弹幕
- [ ] 获取SC  
使用方法：  
```go
b := blivehelper.Default()
err := b.Login(room_id)
if err != nil {
  panic("wrong")
}
ch, _ := b.GetDanmu()
for danmu := range ch {
  if danmu[0] == "2" {
    fmt.Println("礼物：%s", danmu[1])
  } else {
    fmt.Println("弹幕：%s", danmu[1])
  }
}
