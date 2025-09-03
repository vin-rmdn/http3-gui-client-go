# HTTP3 GUI Client in Go
A desktop client application for making HTTP3 requests with a graphical user interface, built with Go and GTK3.

![image](docs/http3%20gui%20icon.png)
(this image is only made within 10 minutes inside Photopea. As much as I love AI I love artists more, so please pay your local artists a visit!)

## Motive
As of September 2025 there are _no_ HTTP GUI client for testing out HTTP3. I created a mock HTTP3 chat server to illustrate improvements in HTTP3, but so far I am blocked on the lack of GUI clients for macOS.

## Features
- HTTP3 protocol support via QUIC
- Intuitive GTK-based user interface
- Support for different HTTP methods (GET, POST)
- Configurable timeout settings
- Structured logging with configurable levels
- Cross-platform compatibility (with focus on macOS)

## Prerequisites
- Go 1.25+
- GTK3 development libraries
- Required Go packages (installed automatically via go.mod):
    - github.com/gotk3/gotk3
    - github.com/quic-go/quic-go
    - gopkg.in/yaml.v3

## Installing GTK3 development libraries
### macOS
```sh
brew install pkg-config gtk+3 adwaita-icon-theme
```

### Ubuntu/Debian
```sh
sudo apt-get install libgtk-3-dev
```

## Configuration
The application uses a YAML configuration file with the following structure:
```yaml
today:
  message: "An exciting opportunity lies ahead of you."

log:
  level: debug

window:
  title: HTTP3 GUI Client in Go

http:
  timeout_in_second: 10
```

Copy config.sample.yml to config.yml and adjust as needed.

## Building
### Building for macOS
This will create a macOS application bundle in the out directory.
```sh
make clean build-macos
```

### Building for other platforms
```sh
go build -o http3-gui-client-go main.go
```

## Usage
1. Launch the application
2. Enter the URL in the text field
3. Select the HTTP method from the dropdown
4. Click "Send request"
5. View the response in the log output

## Project Structure
- main.go: Application entry point
- application: GUI application implementation
- config: Configuration loading and management
- logger: Logging setup and utilities
- macos: macOS specific resources for application bundling

## License
This project is open-source.

## Security Considerations
The application validates input parameters
TLS is properly configured for secure HTTP3 connections
Timeout settings help prevent hanging connections
Error handling is implemented throughout the codebase

## Contributions
Contributions are welcome. Please ensure that code follows the project's style and includes appropriate tests.

## Acknowledgements
Thank you to OpenAI for bringing me up to speed with the GTK library---the documentation could use some rework, and the Go adapter can explain features better.
