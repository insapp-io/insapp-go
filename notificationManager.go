package insapp

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

var firebaseApp = initializeFirebaseApp()

// Please refer to https://firebase.google.com/docs/reference/fcm/rest/v1/projects.messages

func initializeFirebaseApp() *firebase.App {
	firebaseApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	return firebaseApp
}

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

func buildTopicsStringForPromotions(prefix string, promotions []string) string {
	topics := ""

	for _, promotion := range promotions {
		topics += fmt.Sprintf(` || '%s-%s' in topics`, prefix, promotion)
	}

	return topics
}

// TriggerNotificationForUserFromPost will send a Notification as well as a FCM push
// notification to the given user
func TriggerNotificationForUserFromPost(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}
	user := getNotificationUserForUser(receiver)

	sendNotificationToUsers(notification, []NotificationUser{user})

	sendPushNotificationToDevice(GetUser(sender).Username, message, content.Hex(), ".activities.PostActivity", user.Token)
}

// TriggerNotificationForUserFromEvent will send a Notification as well as a FCM push
// notification to the given user
func TriggerNotificationForUserFromEvent(sender bson.ObjectId, receiver bson.ObjectId, content bson.ObjectId, message string, comment Comment, tagType string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Comment: comment, Type: tagType}
	user := getNotificationUserForUser(receiver)

	sendNotificationToUsers(notification, []NotificationUser{user})

	sendPushNotificationToDevice(GetUser(sender).Username, message, content.Hex(), ".activities.EventActivity", user.Token)
}

// TriggerNotificationForEvent will send a Notification as well as a FCM push
// notification to users targeted by the Event platform and promotion
func TriggerNotificationForEvent(event Event, sender bson.ObjectId, content bson.ObjectId, message string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Type: "event"}

	var users []NotificationUser
	var filteredUsers []NotificationUser

	var platforms string

	if contains("iOS", event.Plateforms) && contains("android", event.Plateforms) {
		filteredUsers = getAllUsers()
		platforms = "('events-android' in topics || 'events-ios' in topics)"
	} else if contains("iOS", event.Plateforms) {
		filteredUsers = getiOSUsers("")
		platforms = "'events-ios' in topics"
	} else if contains("android", event.Plateforms) {
		filteredUsers = getAndroidUsers("")
		platforms = "'events-android' in topics"
	}

	for _, notificationUser := range filteredUsers {
		var user = GetUser(notificationUser.UserId)
		if contains(strings.ToUpper(user.Promotion), event.Promotions) {
			users = append(users, notificationUser)
		}
	}

	sendNotificationToUsers(notification, users)

	for _, promotion := range event.Promotions {
		var topics = fmt.Sprintf(`%s && 'events-%s' in topics`, platforms, promotion)
		sendPushNotificationToTopics(event.Name, message, content.Hex(), ".activities.EventActivity", topics)
	}
	sendPushNotificationToTopics(event.Name, message, content.Hex(), ".activities.EventActivity", fmt.Sprintf(`%s && 'events-unknown-class' in topics`, platforms))
}

// TriggerNotificationForPost will send a Notification as well as a FCM push
// notification to users targeted by the Post platform and promotion
func TriggerNotificationForPost(post Post, sender bson.ObjectId, content bson.ObjectId, message string) {
	notification := Notification{Sender: sender, Content: content, Message: message, Type: "post"}

	var users []NotificationUser
	var filteredUsers []NotificationUser

	var platforms string

	if contains("iOS", post.Plateforms) && contains("android", post.Plateforms) {
		filteredUsers = getAllUsers()
		platforms = "('posts-android' in topics || 'posts-ios' in topics)"
	} else if contains("iOS", post.Plateforms) {
		filteredUsers = getiOSUsers("")
		platforms = "'posts-ios' in topics"
	} else if contains("android", post.Plateforms) {
		filteredUsers = getAndroidUsers("")
		platforms = "'posts-android' in topics"
	}

	for _, notificationUser := range filteredUsers {
		var user = GetUser(notificationUser.UserId)
		if contains(strings.ToUpper(user.Promotion), post.Promotions) {
			users = append(users, notificationUser)
		}
	}

	sendNotificationToUsers(notification, users)

	for _, promotion := range post.Promotions {
		var topics = fmt.Sprintf(`%s && 'posts-%s' in topics`, platforms, promotion)
		sendPushNotificationToTopics(post.Title, message, content.Hex(), ".activities.PostActivity", topics)
	}
	sendPushNotificationToTopics(post.Title, message, content.Hex(), ".activities.PostActivity", fmt.Sprintf(`%s && 'posts-unknown-class' in topics`, platforms))
}

func sendNotificationToUsers(notification Notification, users []NotificationUser) {
	for _, user := range users {
		notification.Receiver = user.UserId
		notification = AddNotification(notification)
	}
}

func sendPushNotificationToDevice(title string, message string, objectID string, clickAction string, token string) {
	ctx := context.Background()
	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
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

func sendPushNotificationToTopics(title string, message string, objectID string, clickAction string, topics string) {
	ctx := context.Background()
	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
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
