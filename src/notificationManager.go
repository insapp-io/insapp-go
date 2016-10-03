package main

import (
  apns "github.com/anachronistic/apns"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func getiOSUsers(user bson.ObjectId) []NotificationUser {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification_user")
  var result []NotificationUser
  if user != "" {
      db.Find(bson.M{"os": "iOS"}).All(&result)
  }else{
    db.Find(bson.M{"os": "iOS", "userId": user}).All(&result)
  }
  return result
}

func getOSForUser(user bson.ObjectId) string {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification_user")
  var result NotificationUser
  db.Find(bson.M{"userId": user}).One(&result)
  return result.Os
}

func TriggerNotificationForUser(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment){
  notification := Notification{Sender: sender, Receiver: receiver, Content: content, Message: message, Comment: comment, Type: "tag"}
  if getOSForUser(receiver) == "iOS"{
    triggeriOSNotification(notification)
  }
  // if getOSForUser(receiver) == "Android"{
  //   triggerAndroidNotification(notification)
  // }
}

func TriggerNotificationForEvent(sender bson.ObjectId, content bson.ObjectId, message string){
  notification := Notification{Sender: sender, Content: content, Message: message, Type: "event"}
  triggeriOSNotification(notification)
  //triggerAndroidNotification(notification)
}

func TriggerNotificationForPost(sender bson.ObjectId, content bson.ObjectId, message string){
  notification := Notification{Sender: sender, Content: content, Message: message, Type: "post"}
  triggeriOSNotification(notification)
  //triggerAndroidNotification(notification)
}

func triggeriOSNotification(notification Notification){
  done := make(chan bool)
  users := getiOSUsers(notification.Receiver)
  for _, user := range users {
    notification.Receiver = user.UserId
    AddNotification(notification)
    go sendiOSNotificationToDevice(user.Token, notification, false, done)
  }
  <- done
}

func sendiOSNotificationToDevice(token string, notification Notification, dev bool, done chan bool) {
  payload := apns.NewPayload()
  payload.Alert = notification.Message
  payload.Badge = 42
  payload.Sound = "bingbong.aiff"

  pn := apns.NewPushNotification()
  pn.DeviceToken = token
  pn.AddPayload(payload)
  pn.Set("type", notification.Type)
  pn.Set("sender", notification.Sender)
  pn.Set("content", notification.Content)
  pn.Set("message", notification.Message)

  if dev {
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
