Instructions for running server:

$ cd $GOPATH/src/github.com/patkaehuaea/server
$ go build
$ go install
$ $GOPATH/bin/server

Alternatively:

$ cd $GOPATH/src/github.com/patkaehuaea/server
$ make install
$ $GOPATH/bin/server

To view the version number:

$ $GOPATH/bin/server -V
Version number: v1.0.7

To change default listen port:

$ server -port="9001"

To view main page:

http://localhost:8080/
http://localhost:8080/index.html


To login:

http://localhost:8080/login


To logout:

http://localhost:8080/logout


To view the current time:

$ curl http://localhost:8080/time
<html>
<head>
<style>
p {font-size: xx-large}
span.time {color: red}
</style>
</head>
<body>
<p>The time is now <span class="time">5:59:36 PM</span>.</p>
</body>


To view example 404 page:

$ curl http://localhost:8080/timer
<html>
<body>
<p>These are not the URLs you're looking for.</p>
</body>
</html>

$ curl -I http://localhost:8080/time/is/a/valuable/resource/andshouldnotbewasted
HTTP/1.1 404 Not Found
Date: Wed, 14 Jan 2015 02:00:55 GMT
Content-Length: 80
Content-Type: text/html; charset=utf-8