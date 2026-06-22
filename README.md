# ascii-art-web

## Description

Ascii-art-web is a Go web application that exposes a form-based interface for generating ASCII art in three styles: `standard`, `shadow`, and `thinkertoy`.

## Project structure

- `main.go` — application entry point
- `server/server.go` — HTTP server setup and static file handling
- `handlers/handlers.go` — route registration and request handling
- `ascii-art/AsciiArt.go` — ASCII art generation logic and banner loader
- `ascii-art/artstyles/` — embedded font definitions for each banner style
- `templates/` — HTML layout and page templates
- `static/` — client-side JavaScript and CSS assets
- `tests/main_test.go` — basic unit tests

## Authors

- christianotieno

## Usage

1. Install Go 1.22 or later.
2. Run the application from the repository root:
   ```bash
   go run main.go
   ```
3. Open your browser at `http://localhost:8080`.
4. Enter text, choose a banner style, and submit the form.

## Implementation details: algorithm

The server receives form input on `POST /ascii-art` and converts each supported character into a 6-line ASCII art block.
- The `standard` banner uses a fixed block font.
- The `shadow` banner adds a dotted shadow layer to the standard patterns.
- The `thinkertoy` banner renders the same letter shapes with `@` characters.

Templates are served from the `templates/` directory, while static CSS and JavaScript files are served from `static/`.

## Instructions

- `GET /` renders the main form.
- `POST /ascii-art` accepts text and banner selection, then returns the generated ASCII art.
- HTTP status codes:
  - `200 OK` for successful requests.
  - `400 Bad Request` for invalid inputs.
  - `404 Not Found` when a required resource is missing.
  - `500 Internal Server Error` for unexpected failures.

## Description

Ascii-art-web is a Go web application that converts user input into ASCII art using three banner styles: `standard`, `shadow`, and `thinkertoy`.

## Authors

- christianotieno

## Usage

1. Install Go 1.22 or later.
2. Run the server from the project root:
   ```bash
   go run main.go
   ```
3. Open your browser at `http://localhost:8080`.
4. Enter text, choose a banner style, and submit the form.

## Implementation details: algorithm

The server parses form input from `POST /ascii-art` and maps each supported character to a 6-line ASCII art pattern.
- `standard`: prints block letters with `#` characters.
- `shadow`: renders the same block letters with a dotted shadow to the right.
- `thinkertoy`: renders the block letters as `@` characters.

The result is rendered with Go HTML templates located in the `templates` directory.

## Instructions

- `GET /` serves the main HTML form.
- `POST /ascii-art` accepts text and banner selection and returns the generated ASCII art.
- HTTP status codes are used as follows:
  - `200 OK` for successful requests.
  - `400 Bad Request` for invalid form data.
  - `404 Not Found` when the requested banner or template is unavailable.
  - `500 Internal Server Error` for unexpected server issues.

