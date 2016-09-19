package main

import (
  apns "github.com/anachronistic/apns"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func getiOSUser() []NotificationUser {
  session, _ := mgo.Dial("127.0.0.1")
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  db := session.DB("insapp").C("notification")
  var result []NotificationUser
  db.Find(bson.M{"os": "iOS"}).All(&result)
  return result
}

func getiOSTokenDevice() []string {
  var result []string
  notificationUsers := getiOSUser()
    for _, notif := range notificationUsers {
       result = append(result, notif.Token)
   }
   return result
}

func TriggerNotification(message string){
  triggeriOSNotification(message)
}

func triggeriOSNotification(message string){
  done := make(chan bool)
  devices := getiOSTokenDevice()
  for _, device := range devices {
    go sendiOSNotificationToDevice(device, message, true, done)
  }
  <- done
}

func sendiOSNotificationToDevice(token string, message string, dev bool, done chan bool) {
  payload := apns.NewPayload()
  payload.Alert = message
  payload.Badge = 42
  payload.Sound = "bingbong.aiff"

  pn := apns.NewPushNotification()
  pn.DeviceToken = token
  pn.AddPayload(payload)

  if dev {
    client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "InsappDevCert.pem", "InsappDev.pem")
    client.Send(pn)
    pn.PayloadString()
  }else{
    client := apns.NewClient("gateway.push.apple.com:2195", "InsappDevProd.pem", "InsappProd.pem")
    client.Send(pn)
    pn.PayloadString()
  }

  done <- true
}
