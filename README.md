# socks5-tracker

This library is an extension to [this](https://github.com/things-go/go-socks5) library. It provides some extra features like:

- Limiting global usage
- Limiting usage per user, based on username authentication
- Logging usage (both global & per user) periodically

Library has around `40%` unit test coverage. I intentionally skipped other parts, mostly because there are network connections, where it makes more sense to cover them with integration tests rather then unit tests.

To see how it works, refer to `/examples` folder.