# ascii-art-web

A minimal Go web application that renders ASCII art from banner font files in `banners/`.

## Project structure

- `main.go` — application entry point
- `banner_handler.go` — HTTP handlers and banner-loading logic
- `README.md` — project documentation
- `banners/` — banner font files used to render ASCII art
- `templates/index.html` — form and result page template

## Usage

1. Install Go 1.22 or later.
2. Run the server from the project root:
   ```bash
   go run main.go
   ```
3. Open your browser at `http://localhost:8080`.
4. Enter text and choose a banner style.

## How it works

- The server loads banner definitions from `banners/<name>.txt`.
- `GET /` renders the form.
- `POST /ascii-art` generates ASCII art and redisplays the result on the same page.

## Notes

- The app currently supports `standard`, `shadow`, and `thinkertoy` banner files.
- If a banner file is missing, the app falls back to the default banner list.

