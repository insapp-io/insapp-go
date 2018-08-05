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
| `GET`     | `/posts`                                          | `Get all posts`
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

### Association routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `PUT`     | `/associations/{id}`                              | `Update the association with id {id}`
| `POST`    | `/events`                                         | `Create an event`
| `PUT`     | `/events/{id}`                                    | `Update the event with id {id}`
| `DELETE`  | `/events/{id}`                                    | `Delete the event with id {id}`
| `POST`    | `/posts`                                          | `Create a post`
| `PUT`     | `/posts/{id}`                                     | `Update the post with id {id}`
| `DELETE`  | `/posts/{id}`                                     | `Delete the post with id {id}`
| `POST`    | `/images`                                         | `Post an image`
| `POST`    | `/images/{name}`                                  | `Update the image with name {name}`

### Super user routes

| Type      | Endpoint calls                                    | Description
|-----------|---------------------------------------------------|--------------------------------------
| `POST`    | `/associations`                                   | `Create an association`
| `DELETE`  | `/associations/{id}`                              | `Delete the association with id {id}`
| `GET`     | `/associations/{ownerID}/myassociations`          | `Get the associations owned by the association with id {ownerID}`
| `GET`     | `/users`                                          | `Get all users`