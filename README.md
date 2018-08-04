# insapp-go
Backend part of the Insapp project written in Golang

## Dependencies

Don't forget to install Go dependencies:

```
cd src
go get github.com/gorilla/mux
go get github.com/thomas-bouvier/palette-extractor
go get gopkg.in/mgo.v2
```

## Configuration

Edit the configuration file:

```
cp config.json.dist config.json
vi config.json
```

Attributes `google_email` and `google_password` refer to the credentials of your Google account. `google_key` refers to the Firebase key to be used to send push notifications. `mongo_database_password` refers to the MongoDB password. `env` refers to the environment type and should be set to `dev` or `prod`. Finally, `port` refers to the API port.

## Build & Launch

Check that you have MongoDB running.

```
cd src && go build
```

You can now manually launch your process with `./src`. It is listening on 0.0.0.0:9000 by default.


## API Endpoints

### Public routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `GET`     | `/`                                               | `Index`
| `GET`     | `/credit`                                         | `Get the credits`
| `GET`     | `/legal`                                          | `Get the legal conditions`
| `POST`    | `/login/association`                              | `Log an association in`
| `POST`    | `/login/user`                                     | `Log a user in`
| `POST`    | `/signin/user/{ticket}`                           | `Sign a user in with the ticket {ticket}`

### User routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `GET`     | `/association`                                    | `Get all associations`
| `GET`     | `/association/{id}`                               | `Get the association with id {id}`
| `GET`     | `/event`                                          | `Get all future events`
| `GET`     | `/event/{id}`                                     | `Get the event with id {id}`
| `POST`    | `/event/{id}/participant/{userID}/status/{status}`| `Post the attendee status {status} for the user with id {userID} on the event with id {id}`
| `DELETE`  | `/event/{id}/participant/{userID}`                | `Delete the attendee status of the user with id {userID} on the event with id {id}`
| `POST`    | `/event/{id}/comment`                             | `Post a comment on the event with id {id}`
| `DELETE`  | `/event/{id}/comment/{commentID}`                 | `Delete the comment with id {commentID} on the event with id {id}`
| `GET`     | `/post`                                           | `Get all posts`
| `GET`     | `/post/{id}`                                      | `Get the post with id {id}`
| `POST`    | `/post/{id}/like/{userID}`                        | `Post a like for the user with id {userID} on the post with id {id}`
| `DELETE`  | `/post/{id}/like/{userID}`                        | `Post an un like for the user with id {userID} on the post with id {id}`
| `POST`    | `/post/{id}/comment`                              | `Post a comment on the post with id {id}`
| `DELETE`  | `/post/{id}/comment/{commentID}`                  | `Delete the comment with id {commentID} on the post with id {id}`
| `GET`     | `/user/{id}`                                      | `Get the user with id {id}`
| `PUT`     | `/user/{id}`                                      | `Update the user with id {id}`
| `DELETE`  | `/user/{id}`                                      | `Delete the user with id {id}`
| `POST`    | `/notification`                                   | `Post a notification`
| `GET`     | `/notification/{userID}`                          | `Get all notifications of the user with id {userID}`
| `DELETE`  | `/notification/{userID}/{id}`                     | `Delete the notification with id {id} for the user with id {userID}`
| `PUT`     | `/report/user/{userID}`                           | `Report the user with id {userID}`
| `PUT`     | `/report/{postID}/comment/{commentID}`            | `Report the comment with id {commentID} on the post with id {postID}`
| `POST`    | `/search`                                         | `Search for users, associations, events and posts`
| `POST`    | `/search/users`                                   | `Search for users`
| `POST`    | `/search/associations`                            | `Search for associations`
| `POST`    | `/search/events`                                  | `Search for events`
| `POST`    | `/search/posts`                                   | `Search for posts`

### Association routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `PUT`     | `/association/{id}`                               | `Update the association with id {id}`
| `POST`    | `/event`                                          | `Post an event`
| `PUT`     | `/event/{id}`                                     | `Update the event with id {id}`
| `DELETE`  | `/event/{id}`                                     | `Delete the event with id {id}`
| `POST`    | `/post`                                           | `Post a post`
| `PUT`     | `/post/{id}`                                      | `Update the post with id {id}`
| `DELETE`  | `/post/{id}`                                      | `Delete the post with id {id}`
| `POST`    | `/image`                                          | `Post an image`
| `POST`    | `/image/{name}`                                   | `Update the image with name {name}`

### Super user routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `POST`    | `/association`                                    | `Post an association`
| `DELETE`  | `/association/{id}`                               | `Delete the association with id {id}`
| `GET`     | `/association/{ownerID}/myassociations`           | `Get the associations owned by the association with id {ownerID}`
| `GET`     | `/user`                                           | `Get all users`