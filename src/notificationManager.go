package main

import (
  apns "github.com/anachronistic/apns"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "fmt"
  "strings"
  "net/http"
  "bytes"
)

func getiOSUsers(user string) []NotificationUser {
  conf, _ := Configuration()
  session, _ := mgo.Dial(conf.Database)
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
  conf, _ := Configuration()
  session, _ := mgo.Dial(conf.Database)
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
  conf, _ := Configuration()
  session, _ := mgo.Dial(conf.Database)
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification_user")
  var result NotificationUser
  db.Find(bson.M{"userid": user}).One(&result)
  return result
}

func TriggerNotificationForUser(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string){
  notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}
  user := getNotificationUserForUser(receiver)
  if user.Os == "iOS" {
    triggeriOSNotification(notification, []NotificationUser{user})
  }
  if user.Os == "android" {
    triggerAndroidNotification(notification, []NotificationUser{user})
  }
}

func TriggerNotificationForEvent(event Event, sender bson.ObjectId, content bson.ObjectId, message string){
  notification := Notification{Sender: sender, Content: content, Message: message, Type: "event"}
  iOSUsers := getiOSUsers("")
  users := []NotificationUser{}
  for _, notificationUser := range iOSUsers {
    var user = GetUser(notificationUser.UserId)
    if Contains(strings.ToUpper(user.Promotion), event.Promotions) {
      users = append(users, notificationUser)
    }
  }
  if Contains("iOS", event.Plateforms) {
    triggeriOSNotification(notification, users)
  }
  androidUsers := getAndroidUsers("")
  users = []NotificationUser{}
  for _, notificationUser := range androidUsers {
    var user = GetUser(notificationUser.UserId)
    if Contains(strings.ToUpper(user.Promotion), event.Promotions) {
      users = append(users, notificationUser)
    }
  }
  if Contains("android", event.Plateforms) {
    triggerAndroidNotification(notification, users)
  }
}

func TriggerNotificationForPost(post Post, sender bson.ObjectId, content bson.ObjectId, message string){
  notification := Notification{Sender: sender, Content: content, Message: message, Type: "post"}
  iOSUsers := getiOSUsers("")
  users := []NotificationUser{}
  for _, notificationUser := range iOSUsers {
    var user = GetUser(notificationUser.UserId)
    if Contains(strings.ToUpper(user.Promotion), post.Promotions) {
      users = append(users, notificationUser)
    }
  }
  if Contains("iOS", post.Plateforms) {
    triggeriOSNotification(notification, users)
  }
  androidUsers := getAndroidUsers("")
  users = []NotificationUser{}
  for _, notificationUser := range androidUsers {
    var user = GetUser(notificationUser.UserId)
    if Contains(strings.ToUpper(user.Promotion), post.Promotions) {
      users = append(users, notificationUser)
    }
  }
  if Contains("android", post.Plateforms) {
    triggerAndroidNotification(notification, users)
  }
}

func triggerAndroidNotification(notification Notification, users []NotificationUser){
  done := make(chan bool)
  for _, user := range users {
    notification.Receiver = user.UserId
    notification = AddNotification(notification)
    number := len(GetUnreadNotificationsForUser(user.UserId))
    go sendAndroidNotificationToDevice(user.Token, notification, number, done)
  }
  <- done
}

func triggeriOSNotification(notification Notification, users []NotificationUser){
  done := make(chan bool)
  for _, user := range users {
    notification.Receiver = user.UserId
    notification = AddNotification(notification)
    number := len(GetUnreadNotificationsForUser(user.UserId))
    go sendiOSNotificationToDevice(user.Token, notification, number, done)
  }
  <- done
}

func sendiOSNotificationToDevice(token string, notification Notification, number int, done chan bool) {

  conf, _ := Configuration()
  if conf.Environment != "prod" { return }

  payload := apns.NewPayload()
  payload.Alert = notification.Message
  payload.Badge = number
  payload.Sound = "bingbong.aiff"

  pn := apns.NewPushNotification()
  pn.DeviceToken = token
  pn.AddPayload(payload)
  pn.Set("id", notification.ID)
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

  config, _ := Configuration()
  if config.Environment != "prod" { return }

  url := "https://android.googleapis.com/gcm/send"
  notifJson, _ := json.Marshal(notification)
  var jsonStr = "{\"registration_ids\":[\"" + token + "\"], \"data\":" + string(notifJson) + "}"
  req, _ := http.NewRequest("POST", url, bytes.NewBufferString(jsonStr))



  req.Header.Set("Authorization", "key=" + config.GoogleKey)
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, _ := client.Do(req)

  defer resp.Body.Close()
  fmt.Println("response Status:", resp.Status)

  done <- true
}