package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"

	"firebase.google.com/go/messaging"
)

// Please refer to https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages

func getAllUsers() []NotificationUser {
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("notification_user")

	var result []NotificationUser
	db.Find(bson.M{}).All(&result)

	return result
}

func getiOSUsers(user string) []NotificationUser {
	session := GetMongoSession()
	defer session.Close()
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
	session := GetMongoSession()
	defer session.Close()
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
	session := GetMongoSession()
	defer session.Close()
	db := session.DB("insapp").C("notification_user")

	var result NotificationUser
	db.Find(bson.M{"userid": user}).One(&result)

	return result
}

func buildTopicsStringForPromotions(suffix string, promotions []string) string {
	topics := ""

	for _, promotion := range promotions {
		topics += fmt.Sprintf(` || '%s_%s' in topics`, suffix, promotion)
	}

	return topics
}

func TriggerNotificationForUserFromPost(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}

	user := getNotificationUserForUser(receiver)
	sendNotificationToDevices(
		GetUser(sender).Username,
		message,
		content.Hex(),
		".activities.PostActivity",
		notification,
		[]NotificationUser{user})
}

func TriggerNotificationForUserFromEvent(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}

	user := getNotificationUserForUser(receiver)
	sendNotificationToDevices(
		GetUser(sender).Username,
		message,
		content.Hex(),
		".activities.EventActivity",
		notification,
		[]NotificationUser{user})
}

func TriggerNotificationForEvent(event Event, sender bson.ObjectId, content bson.ObjectId, message string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Type: "event"}

	var users []NotificationUser
	var filteredUsers []NotificationUser

	var platforms string

	if Contains("iOS", event.Plateforms) && Contains("android", event.Plateforms) {
		filteredUsers = getAllUsers()
		platforms = "('events_android' in topics || 'events_ios' in topics)"
	} else if Contains("iOS", event.Plateforms) {
		filteredUsers = getiOSUsers("")
		platforms = "'events_ios' in topics"
	} else if Contains("android", event.Plateforms) {
		filteredUsers = getAndroidUsers("")
		platforms = "'events_android' in topics"
	}

	for _, notificationUser := range filteredUsers {
		var user = GetUser(notificationUser.UserId)
		if Contains(strings.ToUpper(user.Promotion), event.Promotions) {
			users = append(users, notificationUser)
		}
	}

	sendNotificationToTopics(
		event.Name,
		message,
		content.Hex(),
		".activities.EventActivity",
		notification,
		users,
		fmt.Sprintf(
			`%s && ('events_unknown_promotion' in topics %s)`,
			platforms,
			buildTopicsStringForPromotions("events", event.Promotions)))
}

func TriggerNotificationForPost(post Post, sender bson.ObjectId, content bson.ObjectId, message string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Type: "post"}

	var users []NotificationUser
	var filteredUsers []NotificationUser

	var platforms string

	if Contains("iOS", post.Plateforms) && Contains("android", post.Plateforms) {
		filteredUsers = getAllUsers()
		platforms = "('posts_android' in topics || 'posts_ios' in topics)"
	} else if Contains("iOS", post.Plateforms) {
		filteredUsers = getiOSUsers("")
		platforms = "'posts_ios' in topics"
	} else if Contains("android", post.Plateforms) {
		filteredUsers = getAndroidUsers("")
		platforms = "'posts_android' in topics"
	}

	for _, notificationUser := range filteredUsers {
		var user = GetUser(notificationUser.UserId)
		if Contains(strings.ToUpper(user.Promotion), post.Promotions) {
			users = append(users, notificationUser)
		}
	}

	sendNotificationToTopics(
		post.Title,
		message,
		content.Hex(),
		".activities.PostActivity",
		notification,
		users,
		fmt.Sprintf(
			`%s && ('posts_unknown_promotion' in topics %s)`,
			platforms,
			buildTopicsStringForPromotions("posts", post.Promotions)))
}

func sendNotificationToDevices(title string, message string, objectID string, clickAction string, notification Notification, users []NotificationUser) {
	for _, user := range users {
		notification.Receiver = user.UserId
		notification = AddNotification(notification)
		//number := len(GetUnreadNotificationsForUser(user.UserId))
		sendPushNotificationToDevice(title, message, objectID, clickAction, user.Token)
	}
}

func sendPushNotificationToDevice(title string, message string, objectID string, clickAction string, token string) {
	configuration, _ := Configuration()

	ctx := context.Background()
	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	var packageName string
	if configuration.Environment != "prod" {
		packageName = "fr.insapp.insapp.debug"
	} else {
		packageName = "fr.insapp.insapp"
	}

	pushNotification := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Data: map[string]string{
			"ID": objectID,
		},
		Android: &messaging.AndroidConfig{
			RestrictedPackageName: packageName,
			Notification: &messaging.AndroidNotification{
				Sound:       "default",
				Color:       "#ec5d57",
				ClickAction: clickAction,
			},
		},
	}

	// Send a message to the device corresponding to the provided
	// registration token
	response, err := client.Send(ctx, pushNotification)
	if err != nil {
		log.Fatalln(err)
	}

	// Response is a message ID string
	fmt.Println("Successfully sent message:", response)
	fmt.Println("Token:", token)
}

func sendNotificationToTopics(title string, message string, objectID string, clickAction string, notification Notification, users []NotificationUser, topics string) {
	for _, user := range users {
		notification.Receiver = user.UserId
		notification = AddNotification(notification)
		//number := len(GetUnreadNotificationsForUser(user.UserId))
	}

	sendPushNotificationToTopics(title, message, objectID, clickAction, topics)
}

func sendPushNotificationToTopics(title string, message string, objectID string, clickAction string, topics string) {
	configuration, _ := Configuration()

	ctx := context.Background()
	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	var packageName string
	if configuration.Environment != "prod" {
		packageName = "fr.insapp.insapp.debug"
	} else {
		packageName = "fr.insapp.insapp"
	}

	pushNotification := &messaging.Message{
		Condition: topics,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Data: map[string]string{
			"ID": objectID,
		},
		Android: &messaging.AndroidConfig{
			RestrictedPackageName: packageName,
			Notification: &messaging.AndroidNotification{
				Sound:       "default",
				Color:       "#ec5d57",
				ClickAction: clickAction,
			},
		},
	}

	// Send a message to devices subscribed to the combination of topics
	// specified by the provided condition.
	response, err := client.Send(ctx, pushNotification)
	if err != nil {
		log.Fatalln(err)
	}

	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	fmt.Println("Condition:", topics)
}
