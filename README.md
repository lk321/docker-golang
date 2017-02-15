# docker-golang
Docker and golang instructions

### Installation

* Copy project on this route
```sh 
/app/oauth/
```

* Run this commands
```sh
$ cd /app/oauth/
$ chmod +x update.sh 
$ go build
$ docker build -t oauth/prod .
$ ./update.sh oauth
```
