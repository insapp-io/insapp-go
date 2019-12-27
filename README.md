# insapp-go

[![Go Report Card](https://goreportcard.com/badge/github.com/thomas-bouvier/insapp-go)](https://goreportcard.com/report/thomas-bouvier/insapp-go)

Backend part of Insapp written in Go

## Installation

### go get

```bash
go get github.com/thomas-bouvier/insapp-go/...
```

#### API

Check that you have MongoDB running.

```bash
insapp-api
```

The process is listening on 0.0.0.0:9000 by default.

#### CLI

```bash
insapp-cli
```

### Docker

#### Build

```bash
docker build -t insapp-go .
```

#### Run

```bash
docker run -it insapp-go
```

To ssh into this container, use the following command:

```bash
docker exec -it insapp-go ash
```

## Configuration

Edit the configuration file:

```bash
cp config.json.dist config.json
vi config.json
```

Attributes `google_email` and `google_password` refer to the credentials of your Google account. These credentials are used to send emails. `mongo_password` refers to the MongoDB password. `env` refers to the environment type and should be set to `prod`, `dev` or `local`. Finally, `port` refers to the API port.

The FCM HTTP v1 API requires some credentials to send push notifications. The `service-account.json` file can be downloaded from the Firebase Cloud Messaging dashboard, and should be copied at the root of this directory. This way, it will be included in the Docker container.

The authentication mechanism relies on JWT tokens and thus needs a secret key. Please generate the following keys at the root of this directory:

```bash
openssl genrsa -out app.rsa 2048
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

## API Endpoints

### Public routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `GET`     | `/`                                               | `Index`
| `GET`     | `/how-to-post`                                    | `Get the tutorial for posting content`
| `GET`     | `/credit`                                         | `Get the credits`
| `GET`     | `/legal`                                          | `Get the legal conditions`
| `POST`    | `/login/association`                              | `Log an association in`
| `POST`    | `/login/user/{ticket}`                            | `Log a user in with the ticket {ticket} provided by CAS`

### User routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `GET`     | `/associations`                                   | `Get all associations`
| `GET`     | `/associations/{id}`                              | `Get the association with id {id}`
| `GET`     | `/associations/{id}/events`                       | `Get all events of the association with id {id}`
| `GET`     | `/associations/{id}/posts`                        | `Get all posts of the association with id {id}`
| `GET`     | `/events`                                         | `Get all future events`
| `GET`     | `/events/{id}`                                    | `Get the event with id {id}`
| `POST`    | `/events/{id}/attend/{userID}/status/{status}`    | `Post the attendee status {status} for the user with id {userID} on the event with id {id}`
| `DELETE`  | `/events/{id}/attend/{userID}`                    | `Delete the attendee status of the user with id {userID} on the event with id {id}`
| `POST`    | `/events/{id}/comment`                            | `Post a comment on the event with id {id}`
| `DELETE`  | `/events/{id}/comment/{commentID}`                | `Delete the comment with id {commentID} on the event with id {id}`
| `GET`     | `/posts`                                          | `Get all posts. You can provide ?range=[{start},{count}] to get only some posts`
| `GET`     | `/posts/{id}`                                     | `Get the post with id {id}`
| `POST`    | `/posts/{id}/like/{userID}`                       | `Post a like for the user with id {userID} on the post with id {id}`
| `DELETE`  | `/posts/{id}/like/{userID}`                       | `Post an unlike for the user with id {userID} on the post with id {id}`
| `POST`    | `/posts/{id}/comment`                             | `Post a comment on the post with id {id}`
| `DELETE`  | `/posts/{id}/comment/{commentID}`                 | `Delete the comment with id {commentID} on the post with id {id}`
| `GET`     | `/users/{id}`                                     | `Get the user with id {id}`
| `PUT`     | `/users/{id}`                                     | `Update the user with id {id}`
| `DELETE`  | `/users/{id}`                                     | `Delete the user with id {id}`
| `POST`    | `/notifications`                                  | `Create a notification`
| `GET`     | `/notifications/{userID}`                         | `Get all notifications of the user with id {userID}`
| `DELETE`  | `/notifications/{userID}/{id}`                    | `Delete the notification with id {id} for the user with id {userID}`
| `PUT`     | `/report/user/{userID}`                           | `Report the user with id {userID}`
| `PUT`     | `/report/{postID}/comment/{commentID}`            | `Report the comment with id {commentID} on the post with id {postID}`
| `POST`    | `/search`                                         | `Search for users, associations, events and posts`
| `POST`    | `/search/users`                                   | `Search for users`
| `POST`    | `/search/associations`                            | `Search for associations`
| `POST`    | `/search/events`                                  | `Search for events`
| `POST`    | `/search/posts`                                   | `Search for posts`
| `POST`    | `/logout/user`                                    | `Logout the current user`

### Association routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `GET`     | `/association`                                    | `Get the current association`
| `PUT`     | `/associations/{id}`                              | `Update the association with id {id}`
| `POST`    | `/events`                                         | `Create an event`
| `PUT`     | `/events/{id}`                                    | `Update the event with id {id}`
| `DELETE`  | `/events/{id}`                                    | `Delete the event with id {id}`
| `POST`    | `/posts`                                          | `Create a post`
| `PUT`     | `/posts/{id}`                                     | `Update the post with id {id}`
| `DELETE`  | `/posts/{id}`                                     | `Delete the post with id {id}`
| `POST`    | `/images`                                         | `Post an image`
| `POST`    | `/logout/association`                             | `Logout the current association user`

### Super user routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `POST`    | `/associations`                                   | `Create an association`
| `DELETE`  | `/associations/{id}`                              | `Delete the association with id {id}`
| `GET`     | `/associations/{ownerID}/myassociations`          | `Get the associations owned by the association with id {ownerID}`
| `GET`     | `/users`                                          | `Get all users`
