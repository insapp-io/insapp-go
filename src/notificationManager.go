package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/anachronistic/apns"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

// fcmResponseStatus represents fcm response message
type fcmResponseStatus struct {
	Ok            bool
	StatusCode    int
	MulticastId   int64               `json:"multicast_id"`
	Success       int                 `json:"success"`
	Fail          int                 `json:"failure"`
	Canonical_ids int                 `json:"canonical_ids"`
	Results       []map[string]string `json:"results,omitempty"`
	MsgId         int64               `json:"message_id,omitempty"`
	Err           string              `json:"error,omitempty"`
	RetryAfter    string
}

func getiOSUsers(user string) []NotificationUser {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("notification_user")

	var result []NotificationUser
	if user == "" {
		db.Find(bson.M{"os": "iOS"}).All(&result)
	} else {
		db.Find(bson.M{"os": "iOS", "userid": user}).All(&result)
	}

	return result
}

func getAndroidUsers(user string) []NotificationUser {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("notification_user")

	var result []NotificationUser
	if user == "" {
		db.Find(bson.M{"os": "android"}).All(&result)
	} else {
		db.Find(bson.M{"os": "android", "userid": user}).All(&result)
	}

	return result
}

func getNotificationUserForUser(user bson.ObjectId) NotificationUser {
	_, info, _ := Configuration()
	session, _ := mgo.DialWithInfo(info)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("insapp").C("notification_user")

	var result NotificationUser
	db.Find(bson.M{"userid": user}).One(&result)

	return result
}

func TriggerNotificationForUserFromPost(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}

	user := getNotificationUserForUser(receiver)
	if user.Os == "iOS" {
		triggeriOSNotification(notification, []NotificationUser{user})
	}
	if user.Os == "android" {
		triggerAndroidNotification(GetUser(sender).Username, message, content.Hex(), ".activities.PostActivity", notification, []NotificationUser{user})
	}
}

func TriggerNotificationForUserFromEvent(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}

	user := getNotificationUserForUser(receiver)
	if user.Os == "iOS" {
		triggeriOSNotification(notification, []NotificationUser{user})
	}
	if user.Os == "android" {
		triggerAndroidNotification(GetUser(sender).Username, message, content.Hex(), ".activities.EventActivity", notification, []NotificationUser{user})
	}
}

func TriggerNotificationForEvent(event Event, sender bson.ObjectId, content bson.ObjectId, message string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Type: "event"}
	var users []NotificationUser

	if Contains("iOS", event.Plateforms) {
		iOSUsers := getiOSUsers("")
		for _, notificationUser := range iOSUsers {
			var user = GetUser(notificationUser.UserId)
			if Contains(strings.ToUpper(user.Promotion), event.Promotions) {
				users = append(users, notificationUser)
			}
		}

		triggeriOSNotification(notification, users)
	}

	if Contains("android", event.Plateforms) {
		androidUsers := getAndroidUsers("")
		users = []NotificationUser{}
		for _, notificationUser := range androidUsers {
			var user = GetUser(notificationUser.UserId)
			if Contains(strings.ToUpper(user.Promotion), event.Promotions) {
				users = append(users, notificationUser)
			}
		}
		for _, user := range users {
			notification.Receiver = user.UserId
			AddNotification(notification)
		}

		sendAndroidNotificationToTopics([]string{"events"}, event.Name, message, content.Hex(), ".activities.EventActivity")
	}
}

func TriggerNotificationForPost(post Post, sender bson.ObjectId, content bson.ObjectId, message string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Type: "post"}
	var users []NotificationUser

	if Contains("iOS", post.Plateforms) {
		iOSUsers := getiOSUsers("")
		for _, notificationUser := range iOSUsers {
			var user = GetUser(notificationUser.UserId)
			if Contains(strings.ToUpper(user.Promotion), post.Promotions) {
				users = append(users, notificationUser)
			}
		}

		triggeriOSNotification(notification, users)
	}

	if Contains("android", post.Plateforms) {
		androidUsers := getAndroidUsers("")
		users = []NotificationUser{}
		for _, notificationUser := range androidUsers {
			var user = GetUser(notificationUser.UserId)
			if Contains(strings.ToUpper(user.Promotion), post.Promotions) {
				users = append(users, notificationUser)
			}
		}
		for _, user := range users {
			notification.Receiver = user.UserId
			AddNotification(notification)
		}

		sendAndroidNotificationToTopics([]string{"news"}, post.Title, message, content.Hex(), ".activities.PostActivity")
	}
}

func triggerAndroidNotification(title string, message string, objectId string, clickAction string, notification Notification, users []NotificationUser) {
	for _, user := range users {
		notification.Receiver = user.UserId
		notification = AddNotification(notification)
		//number := len(GetUnreadNotificationsForUser(user.UserId))
		sendAndroidNotificationToDevice(user.Token, title, message, objectId, clickAction)
	}
}

func triggeriOSNotification(notification Notification, users []NotificationUser) {
	for _, user := range users {
		notification.Receiver = user.UserId
		notification = AddNotification(notification)
		number := len(GetUnreadNotificationsForUser(user.UserId))
		sendiOSNotificationToDevice(user.Token, notification, number)
	}
}

func sendiOSNotificationToDevice(token string, notification Notification, number int) {
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

	configuration, _, _ := Configuration()

	if configuration.Environment != "prod" {
		client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "InsappDevCert.pem", "InsappDev.pem")
		client.Send(pn)
		pn.PayloadString()
	} else {
		client := apns.NewClient("gateway.push.apple.com:2195", "InsappProdCert.pem", "InsappProd.pem")
		client.Send(pn)
		pn.PayloadString()
	}
}

func sendAndroidNotificationToDevice(token string, title string, message string, objectId string, clickAction string) {
	url := "https://fcm.googleapis.com/fcm/send"

	var jsonStr string
	configuration, _, _ := Configuration()

	if configuration.Environment != "prod" {
		jsonStr = fmt.Sprintf(`{
			"to":"%s",
			"notification":{"title":"%s","body":"%s","sound":"default","color":"#ec5d57","click_action":"%s"},
			"data":{"ID":"%s"},
			"restricted_package_name":"fr.insapp.insapp.debug"
			}`, token, title, message, clickAction, objectId)
	} else {
		jsonStr = fmt.Sprintf(`{
			"to":"%s",
			"notification":{"title":"%s","body":"%s","sound":"default","color":"#ec5d57","click_action":"%s"},
			"data":{"ID":"%s"},
			"restricted_package_name":"fr.insapp.insapp"
			}`, token, title, message, clickAction, objectId)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+configuration.FirebaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("Android notification response :")
	fmt.Println("Token:", token)
	fmt.Println("Status:", resp.StatusCode)

	var res fcmResponseStatus

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &res)

	spew.Dump(res)
}

func sendAndroidNotificationToTopics(topics []string, title string, message string, objectId string, clickAction string) {
	url := "https://fcm.googleapis.com/fcm/send"

	var jsonStr string
	configuration, _, _ := Configuration()

	var topicsStr string

	if topics == nil || len(topics) == 0 {
		return
	} else if len(topics) == 1 {
		topicsStr = "/topics/" + topics[0]

		if configuration.Environment != "prod" {
			jsonStr = fmt.Sprintf(`{
				"to":"%s",
				"notification":{"title":"%s","body":"%s","sound":"default","color":"#ec5d57","click_action":"%s"},
				"data":{"ID":"%s"},
				"restricted_package_name":"fr.insapp.insapp.debug"
				}`, topicsStr, title, message, clickAction, objectId)
		} else {
			jsonStr = fmt.Sprintf(`{
				"to":"%s",
				"notification":{"title":"%s","body":"%s","sound":"default","color":"#ec5d57","click_action":"%s"},
				"data":{"ID":"%s"},
				"restricted_package_name":"fr.insapp.insapp"
				}`, topicsStr, title, message, clickAction, objectId)
		}
	} else {
		for i := 0; i < len(topics); i++ {
			if i > 0 {
				topicsStr += " || "
			}
			topicsStr += "'" + topics[i] + "' in topics"
		}

		if configuration.Environment != "prod" {
			jsonStr = fmt.Sprintf(`{
				"condition":"%s",
				"notification":{"title":"%s","body":"%s","sound":"default","color":"#ec5d57","click_action":"%s"},
				"data":{"ID":"%s"},
				"restricted_package_name":"fr.insapp.insapp.debug"
				}`, topicsStr, title, message, clickAction, objectId)
		} else {
			jsonStr = fmt.Sprintf(`{
				"condition":"%s",
				"notification":{"title":"%s","body":"%s","sound":"default","color":"#ec5d57","click_action":"%s"},
				"data":{"ID":"%s"},
				"restricted_package_name":"fr.insapp.insapp"
				}`, topicsStr, title, message, clickAction, objectId)
		}
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+configuration.FirebaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("Android notification response:")
	fmt.Println("Condition:", topicsStr)
	fmt.Println("Status:", resp.StatusCode)

	var res fcmResponseStatus

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &res)

	spew.Dump(res)
}
