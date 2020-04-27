# brevis [![Go Report Card](https://goreportcard.com/badge/github.com/admiralobvious/brevis)](https://goreportcard.com/report/github.com/admiralobvious/brevis)

brevis is very simple URL shortener built with [echo](https://github.com/labstack/echo) and using MongoDB for its database.

### Building & Running locally

```shell script
go build ./cmd/brevis && ./brevis
⇛ http server started on 127.0.0.1:1323
```

### Usage
```shell script
Usage of ./brevis:
      --app-name string                     The name of the application. Used to prefix environment variables. (default "brevis")
      --base-url string                     Base URL to prefix short URLs with (default "http://localhost:1323/")
      --bind-address ip                     The IP address to listen at. (default 127.0.0.1)
      --bind-port uint                      The port to listen at. (default 1323)
      --cors-allow-credentials              Tells browsers whether to expose the response to frontend JavaScript code when the request's credentials mode (Request.credentials) is 'include'.
      --cors-allow-headers strings          Indicate which HTTP headers can be used during an actual request.
      --cors-allow-methods strings          Indicates which HTTP methods are allowed for cross-origin requests. (default [GET,HEAD,PUT,PATCH,POST,DELETE])
      --cors-allow-origins strings          Indicates whether the response can be shared with requesting code from the given origin (default [*])
      --cors-expose-headers strings         Indicates which headers can be exposed as part of the response by listing their name.
      --cors-max-age int                    Indicates how long the results of a preflight request can be cached.
      --database-mongodb-password string    MongoDB password
      --database-mongodb-timeout duration   Timeout connecting/reading/writing to MongoDB (default 5s)
      --database-mongodb-uri string         URI of the MongoDB server (default "mongodb://127.0.0.1")
      --database-mongodb-username string    MongoDB username
      --database-type string                Type of database to use to store short URLs (default "mongodb")
      --env-name string                     The environment of the application. Used to load the right config file. (default "local")
      --graceful-timeout uint               Timeout for graceful shutdown. (default 30)
      --log-file string                     The log file to write to. 'stdout' means log to stdout, 'stderr' means log to stderr and 'null' means discard log messages. (default "stdout")
      --log-format string                   The log format. Valid format values are: text, json. (default "text")
      --log-level string                    The granularity of log outputs. Valid log levels: debug, info, warning, error and critical. (default "info")
      --log-requests-disabled               Disables HTTP requests logging.
```

### API:

Shorten a long URL:
```shell script
curl -H "Content-Type: application/json" -X POST -d '{"url":"https://google.com"}' http://localhost:1323/shorten
{"short_url":"http://localhost:1323/arQZNnaKt"}
```

Unshorten a short URL:
```shell script
curl -H "Content-Type: application/json" -X POST -d '{"short_url":"arQZNnaKt"}' http://localhost:1323/unshorten
{"url":"https://google.com"}
```

Get redirected to the long URL (for browsers mostly):
```shell script
curl -i http://localhost:1323/arQZNnaKt
HTTP/1.1 301 Moved Permanently
Access-Control-Allow-Origin: *
Location: https://google.com
Vary: Origin
Date: Wed, 11 Jan 2017 20:26:35 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```

Get stats:
```shell script
curl http://localhost:1323/arQZNnaKt/stats
{
    "created_at": "2020-04-27T19:40:56.495Z",
    "last_accessed_at": "2020-04-27T19:47:50.709Z",
    "last_updated_at": "2020-04-27T19:47:50.712Z",
    "referrers": [
        {
            "address": "",
            "first_visit_at": "2020-04-27T19:47:50.709Z",
            "last_visit_at": "2020-04-27T19:47:50.709Z",
            "visits": 1
        }
    ],
    "short_url": "http://localhost:1323/arQZNnaKt",
    "unique_views": 1,
    "url": "https://google.com",
    "views": 1
}

```

### Building Docker image
```shell script
docker build -t brevis .
```

### Using brevis on your own server:
1. Modify `config.prod.toml` in `configs/` to add at least `bind-url` with your own domain e.g. `https://u.mydomain.com/`.
You will probably need to add the relevant settings for MongoDB as well.
1. Build your own docker image.
1. Deploy!

### Kubernetes
If you're using Kubernetes, you can use the exising [image](https://hub.docker.com/repository/docker/admiralobvious/brevis) and you just have to set the environment variables in your manifest:
```yaml
env:
- name: BREVIS_BASE_URL
  value: https://u.mydomain.com/
```
