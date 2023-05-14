# socks5-tracker

This library is an extension to [this](https://github.com/things-go/go-socks5) library. It provides some extra features like:

- Limiting global usage
- Limiting usage per user, based on username authentication
- Logging usage (both global & per user) periodically

Library has around `40%` unit test coverage. You can check it by running `go test -cover ./...` command in root folder. 
I intentionally skipped other parts, mostly because there are network connections, where it makes more sense to cover them with integration tests rather then unit tests.

To see how it works, refer to `/examples` folder. You can use this command to test it after getting it run with `go run .`:

`curl --socks5 giray:pass@localhost:8000 https://www.youtube.com`

There is also another http endpoint at `localhost:8081/logs` where external admins can check how much traffic has been used both globally and per user basis.