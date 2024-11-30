A `config/atc.yml` file contains two sections: `server` and `routes`.

### server

This section is completely optional.  Every field listed below has a default, and can be overriddenby setting an environment variable with the same name, just in upper case.

| Field Name            | Description                                             | Default Value |
|-----------------------|---------------------------------------------------------|---------------|
| `target_port`         | The port that your Puma server should run on. Thruster will set `PORT` to this value when starting your server. | 3000 |
| `cache_size`          | The size of the HTTP cache in bytes. | 64MB |
| `max_cache_item_size` | The maximum size of a single item in the HTTP cache in bytes. | 1MB |
| `x_sendfile_enabled`  | Whether to enable X-Sendfile support. Set to `0` or `false` to disable. | Enabled |
| `max_request_body`    | The maximum size of a request body in bytes. Requests larger than this size will be refused; `0` means no maximum size is enforced. | `0` |
| `bad_gateway_page`    | Path to an HTML file to serve when the backend server returns a 502 Bad Gateway error. If there is no file at the specific path, Thruster will serve an empty 502 response instead. | `./public/502.html` |
| `http_port`           | The port to listen on for HTTP traffic. | 3000 (development)<br>8080 (production) |
| `http_idle_timeout`   | The maximum time in seconds that a client can be idle before the connection is closed. | 60 |
| `http_read_timeout`   | The maximum time in seconds that a client can take to send the request headers and body. | 30 |
| `http_write_timeout`  | The maximum time in seconds during which the client must read the response. | 30 |
| `forward_headers`     | Whether to forward X-Forwarded-* headers from the client. | Disabled when running with TLS; enabled otherwise |
| `debug`               | Set to `1` or `true` to enable debug logging. | Disabled |
| `health_check_path`   | Path used to check if the app is ready to accept requests | `/up` |

### routes

An array of routes to tenant applications.  Each contains:

| Field Name            | Description                                             | Default Value |
|-----------------------|---------------------------------------------------------|---------------|
| `name`                | Name of route.  Used to create PID file and passed to the application in the `FLY_ATC_NAME` environment variable. |
| `endpoint`            | Prefix path.  Used to route requests and passed to the application in the `FLY_ATC_SCOPE` environment variable. |
| `database`            | Name of the database.  `.sqlite3` is appended to this name, and the result is used to replace the basename specified in `DATABASE_URL`. |
| `region`              | Region where this application is hosted. |
| `instance`            | Machine id where this application is hosted.  Takes precedence over `region` |

If no routes are specified, a single route with a blank name and endpoint is created.

