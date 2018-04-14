# insapp-go
Backend part of the Insapp project written in Go-lang

## Dependencies

Don't forget to install Go dependencies:

```
cd src
go get github.com/gorilla/mux
go get gopkg.in/mgo.v2
```

## Configuration

Edit the configuration file:

```
cp config.json.dist config.json
vi config.json
```

Attributes `email` and `password` refer to the credentials of your Google account. `googlekey` refers to the Firebase key to be used to send push notifications.

## Build & Launch

Check that you have mongodb running

```
cd src && go build
```

You can now manually launch your process with `./src`. It is listening on 0.0.0.0:9000 by default.


## API Endpoints

### Association

@TODO

### Event

@TODO

### User

@TODO

### Post

@TODO
