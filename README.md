# brevis

## How to Use

Run brevis:
```
go build && ./brevis
⇛ http server started on 0.0.0.0:1323
```

Shorten a long URL:
```
$ curl -H "Content-Type: application/json" -X POST -d '{"url":"http://google.com"}' http://localhost:1323/shorten
{"url":"http://google.com","shortUrl":"http://localhost:1323/onZB1cMZj"}
```

Unshorten a short URL:
```
$ curl -H "Content-Type: application/json" -X POST -d '{"shortUrl":"onZB1cMZj"}' http://localhost:1323/unshorten
{"url":"http://google.com","shortUrl":"http://localhost:1323/onZB1cMZj"}
```

Get redirected to the long URL (for browsers mostly):
```
$ curl -i http://localhost:1323/onZB1cMZj
HTTP/1.1 301 Moved Permanently
Access-Control-Allow-Origin: *
Location: http://google.com
Vary: Origin
Date: Wed, 11 Jan 2017 20:26:35 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```
