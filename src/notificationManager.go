package main

import (
  apns "github.com/anachronistic/apns"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "fmt"
  "net/http"
  "bytes"
)

func getiOSUsers(user string) []NotificationUser {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification_user")
  var result []NotificationUser
  if user == "" {
    db.Find(bson.M{"os": "iOS"}).All(&result)
  }else{
    db.Find(bson.M{"os": "iOS", "userid": user}).All(&result)
  }
  return result
}

func getAndroidUsers(user string) []NotificationUser {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification_user")
  var result []NotificationUser
  if user == "" {
    db.Find(bson.M{"os": "android"}).All(&result)
  }else{
    db.Find(bson.M{"os": "android", "userid": user}).All(&result)
  }
  return result
}

func getNotificationUserForUser(user bson.ObjectId) NotificationUser {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification_user")
  var result NotificationUser
  db.Find(bson.M{"userid": user}).One(&result)
  return result
}

func TriggerNotificationForUser(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment){
  notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: "tag"}
  user := getNotificationUserForUser(receiver)
  if user.Os == "iOS" {
    triggeriOSNotification(notification, []NotificationUser{user})
  }
  if user.Os == "android" {
    triggerAndroidNotification(notification, []NotificationUser{user})
  }
}

func TriggerNotificationForEvent(sender bson.ObjectId, content bson.ObjectId, message string){
  notification := Notification{Sender: sender, Content: content, Message: message, Type: "event"}
  iOSUsers := getiOSUsers("")
  androidUsers := getAndroidUsers("")
  triggeriOSNotification(notification, iOSUsers)
  triggerAndroidNotification(notification, androidUsers)
}

func TriggerNotificationForPost(sender bson.ObjectId, content bson.ObjectId, message string){
  notification := Notification{Sender: sender, Content: content, Message: message, Type: "post"}
  iOSUsers := getiOSUsers("")
  androidUsers := getAndroidUsers("")
  triggeriOSNotification(notification, iOSUsers)
  triggerAndroidNotification(notification, androidUsers)
}

func triggerAndroidNotification(notification Notification, users []NotificationUser){
  done := make(chan bool)
  for _, user := range users {
    notification.Receiver = user.UserId
    AddNotification(notification)
    number := len(GetUnreadNotificationsForUser(user.UserId))
    go sendAndroidNotificationToDevice(user.Token, notification, number, done)
  }
  <- done
}

func triggeriOSNotification(notification Notification, users []NotificationUser){
  done := make(chan bool)
  for _, user := range users {
    notification.Receiver = user.UserId
    AddNotification(notification)
    number := len(GetUnreadNotificationsForUser(user.UserId))
    go sendiOSNotificationToDevice(user.Token, notification, number, done)
  }
  <- done
}

func sendiOSNotificationToDevice(token string, notification Notification, number int, done chan bool) {

  payload := apns.NewPayload()
  payload.Alert = notification.Message
  payload.Badge = number
  payload.Sound = "bingbong.aiff"

  pn := apns.NewPushNotification()
  pn.DeviceToken = token
  pn.AddPayload(payload)
  pn.Set("type", notification.Type)
  pn.Set("sender", notification.Sender)
  pn.Set("content", notification.Content)
  pn.Set("message", notification.Message)
  if notification.Type == "tag" {
    pn.Set("comment", notification.Comment.ID)
  }

  config, _ := Configuration()

  if config.Environment == "staging" {
    client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "InsappDevCert.pem", "InsappDev.pem")
    client.Send(pn)
    pn.PayloadString()
  }else{
    client := apns.NewClient("gateway.push.apple.com:2195", "InsappProdCert.pem", "InsappProd.pem")
    client.Send(pn)
    pn.PayloadString()
  }

  done <- true
}

func sendAndroidNotificationToDevice(token string, notification Notification, number int, done chan bool) {
  url := "https://android.googleapis.com/gcm/send"
  notifJson, _ := json.Marshal(notification)
  var jsonStr = "{\"registration_ids\":[\"" + token + "\"], \"data\":" + string(notifJson) + "}"
  req, err := http.NewRequest("POST", url, bytes.NewBufferString(jsonStr))

  config, _ := Configuration()

  req.Header.Set("Authorization", "key=" + config.GoogleKey)
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()
  fmt.Println("response Status:", resp.Status)

  done <- true
}
