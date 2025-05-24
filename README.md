# xfwder
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

xfwder is a tool that helps with passing data via custom URLs in macOS. It acts as a bridge between custom URL schemes and applications by forwarding data received through custom URLs to applications via UNIX domain sockets.

## How It Works

The xfwder tool translates custom URL requests into HTTP requests over UNIX domain sockets:

- A request to `xfwder://uds_file/path?queries`
- Gets forwarded as a POST request to `http+unix:/tmp/uds_file.sock/path?queries`

## License
This application is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
