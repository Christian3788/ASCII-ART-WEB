# ASCII Art Web

ASCII Art Web is a simple Go web application that converts user input text into ASCII art using selectable font styles. It provides a clean browser interface for choosing a style, entering text, and rendering the generated art.

## Features

- Generate ASCII art from text input
- Choose from three built-in styles: `standard`, `shadow`, and `thinkertoy`
- Supports multiline input via `\n`
- Simple web interface with form-based submission
- Static assets served from `static/` and HTML templates in `templates/`

## Project Structure

- `main.go` — application entry point
- `server/server.go` — HTTP server setup and static file handling
- `handlers/handlers.go` — route registration and request handling
- `ascii-art/AsciiArt.go` — ASCII art generation logic and embedded font files
- `ascii-art/artstyles/` — embedded font definitions
- `templates/` — HTML layout and page templates
- `static/` — client-side JavaScript and CSS
- `tests/main_test.go` — basic unit tests

## Requirements

- Go 1.21+

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/ASCII-ART-WEB.git
   cd ASCII-ART-WEB
   ```
2. Ensure Go is installed and configured on your machine.

## Running the Application

Start the app with:

```bash
go run main.go
```

Then open your browser at:

```text
http://localhost:8080
```

## Usage

1. Open the app in your browser.
2. Enter the text you want to convert into ASCII art.
3. Select one of the available styles.
4. Click **Generate**.
5. The generated ASCII art will appear below the form.

## Available Styles

- `standard`
- `shadow`
- `thinkertoy`

## Testing

Run the unit tests with:

```bash
go test ./tests
```

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests for bug fixes, improvements, or new features.

## License

This project is licensed under the MIT License.
