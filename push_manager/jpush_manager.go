package push_manager

import (
  "fmt"
  "github.com/ylywyn/jpush-api-go-client"
)

const (
    appKey = "1c080640aa9e7a1462feefc0"
    secret = "7dab3bb63f447b204a92f385"
)

func PushJPushMessage(registration_id, mapkey string) bool {

  //Platform
  var pf jpushclient.Platform
  pf.Add(jpushclient.ANDROID)

  //Audience
  var ad jpushclient.Audience
  ad.SetID([]string{registration_id})

  //Message
  var msg jpushclient.Message
  msg.Title = "map_key"
  msg.Content = mapkey

  payload := jpushclient.NewPushPayLoad()
  payload.SetPlatform(&pf)
  payload.SetAudience(&ad)
  payload.SetMessage(&msg)

  //push
  b, _ := payload.ToBytes()
  c := jpushclient.NewPushClient(secret, appKey)
  str, err := c.Send(b)
  if err != nil {
    fmt.Printf("err:%s\n", err.Error())
    return false
  } else {
    fmt.Printf("ok:%s\n", str)
    return true
  }
}
