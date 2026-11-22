# ascii-art-web

A browser-based ASCII art generator built with Go and packaged with Docker.

---

## Description

`ascii-art-web` is a web application that brings the ASCII art generator to the browser. Instead of running a CLI command, users type text, choose a banner style, and receive the rendered ASCII art directly on the page.

The server is written in Go using only the standard library. It can be run locally with `go run .` or as a Docker container built from the included `Dockerfile`.

Supported banner styles:
- `standard`
- `shadow`
- `thinkertoy`

---

## Features

- Web interface with a text area and a banner selector
- Three ASCII banner styles
- Multi-line input, using either the Enter key or a literal `\n`
- Strict input validation with appropriate HTTP status codes
- Responsive dark theme styled with CSS custom properties
- Multi-stage Docker build with image metadata
- No external dependencies — standard Go library only

---

## Authors

- **Omar Rady** — [@orady](https://01.nextera.education/git/orady)
- **Hazem Elmahdy** — [@helmahdy](https://01.nextera.education/git/helmahdy)

---

## Usage

### Run locally

**1. Clone the repository and enter it:**

```bash
git clone https://01.nextera.education/git/orady/ascii-art-web
cd ascii-art-web
```

**2. Start the server:**

```bash
go run .
```

**3. Open the application in your browser:**

```
http://localhost:8080
```

**4. Enter your text, select a banner, and click Generate.**

### Run with Docker

**1. Build the image:**

```bash
docker build -t ascii-art-web .
```

**2. Start a container:**

```bash
docker run -d --name ascii-art-web -p 8080:8080 ascii-art-web
```

**3. Open `http://localhost:8080` in your browser.**

**4. Stop and clean up when finished:**

```bash
docker stop ascii-art-web
docker rm ascii-art-web
docker image prune -f      # remove dangling images
docker builder prune -f    # remove unused build cache
```

---

## Project Structure

```
ascii-art-web/
├── main.go                    # Entry point — registers routes and starts the server
├── go.mod
├── Dockerfile                 # Two-stage build (builder + runtime)
├── handler/
│   └── handler.go             # HTTP handlers: HomePageHandler and AsciiArtHandler
├── banner/
│   └── banner.go              # Loads banner files and builds the character map
├── render/
│   └── render.go              # Renders one line of text as 8-line ASCII art
├── renderpage/
│   └── render.go              # Parses the template once and renders the page
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt         # Character definitions (96 printable characters each)
├── templates/
│   └── index.html             # HTML page served on GET /
└── static/
    └── style.css              # Page styling
```

---

## System Architecture

Every request enters through `main.go`, is validated and orchestrated by the handler, and leaves as an HTML page rendered from a single template.

```
Browser ──── HTTP request ────▶ main.go (router)
                                     │
                                     ▼
                             handler/handler.go
                              validates the input
                    ┌────────────────┴────────────────┐
                    ▼                                 ▼
            banner/banner.go                  render/render.go
         banners/*.txt → charMap          charMap + text → ASCII lines
                    └────────────────┬────────────────┘
                                     ▼
                           renderpage/render.go
                     templates/index.html → HTML page
                                     │
Browser ◀─── HTTP response ──────────┘
            (status code + HTML)
```

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Routing | `main` | Registers `/`, `/ascii-art`, `/static/` and starts the server |
| HTTP logic | `handler` | Checks method and input, picks the banner, assembles the result |
| Data | `banner` | Reads a banner file and builds a `map[rune][]string` |
| Rendering logic | `render` | Converts one line of text into 8 lines of ASCII art |
| Presentation | `renderpage` | Executes the HTML template and writes the response |

---

## Implementation Details

**Server startup**

`main()` registers three routes and listens on port `8080`:
- `/` → `handler.HomePageHandler`
- `/ascii-art` → `handler.AsciiArtHandler`
- `/static/` → file server for the CSS

The HTML template is parsed **once**, in the `init()` function of the `renderpage` package, so requests only execute it and never re-read the file from disk. If the template is missing or invalid, the program stops at startup instead of failing on every request.

---

**GET /**

`HomePageHandler` returns `404` for any path other than `/`, returns `405` for any method other than `GET`, then renders the empty form.

---

**POST /ascii-art**

`AsciiArtHandler` processes the submission in this order:

1. Rejects any method other than `POST` with `405 Method Not Allowed`
2. Reads the `text` and `banner` form values
3. Normalises line endings (`\r\n` and `\r` are removed or converted, a literal `\n` becomes a newline)
4. Rejects empty or whitespace-only text with `400 Bad Request`
5. Rejects any character outside printable ASCII (32–126), except newlines, with `400`
6. Accepts only `standard`, `shadow` or `thinkertoy` as the banner — anything else returns `400`
7. Loads `banners/<name>.txt` — returns `404` if the file cannot be read
8. Checks that the banner file has at least 855 lines — returns `500` if it is corrupted
9. Builds the `map[rune][]string` that maps each character to its 8-line drawing
10. Splits the text on newlines and renders each line with `render.RenderText`; an empty line produces a blank line
11. Passes the final string to `renderpage.RenderPage`

---

**Template rendering**

`RenderPage` executes the template into an in-memory `bytes.Buffer` first and writes it to the response only if execution succeeded. A failing template therefore returns a clean `500` and never sends a half-written page.

The template shows the result only when one exists:

```html
{{if .Result}}
<pre>{{.Result}}</pre>
{{end}}
```

`html/template` escapes the output automatically.

---

**HTTP status codes**

| Code | When it is returned |
|------|---------------------|
| `200 OK` | Successful `GET /` or `POST /ascii-art` |
| `400 Bad Request` | Empty text, unsupported characters, or an invalid banner name |
| `404 Not Found` | Unknown path, or banner file not found |
| `405 Method Not Allowed` | Wrong method on `/` or on `/ascii-art` |
| `500 Internal Server Error` | Template execution failure, or a corrupted banner file |

---

## Error Handling

| Situation | Response |
|-----------|----------|
| Request to an unknown path | `404 Not Found` |
| `POST` request to `/` | `405 Method Not Allowed` |
| `GET` request to `/ascii-art` | `405 Method Not Allowed` |
| Empty or whitespace-only text | `400 Bad Request` |
| Characters outside ASCII 32–126 | `400 Bad Request` |
| Banner name missing or not one of the three supported | `400 Bad Request` |
| Banner file does not exist on disk | `404 Not Found` |
| Banner file has fewer than 855 lines | `500 Internal Server Error` |
| Template fails while rendering | `500 Internal Server Error` |

---

## Docker

### Multi-stage build

The `Dockerfile` has two stages, so the final image contains only what is needed to run the server:

```
Stage 1 · builder   golang:1.25.6-alpine
          go.mod → go mod download → source code → go build
          output: one static binary (ascii-art-web)

Stage 2 · runtime   alpine:latest
          copies only: the binary + templates/ + static/ + banners/
          exposes 8080 and runs ./ascii-art-web
```

The Go toolchain, the source code and the build cache stay in the first stage and are not part of the final image.

### Dockerfile practices applied

- **Multi-stage build** — separate build and run stages
- **Layer caching** — `go.mod` is copied and dependencies downloaded before the source code, so source changes do not repeat that step
- **Static binary** — `CGO_ENABLED=0` produces a binary that runs on minimal Alpine
- **Small runtime image** — Alpine, with only the files the server reads at runtime
- **Pinned toolchain** — the builder image uses the same Go version as `go.mod`

### Metadata

The final image carries labels describing it:

```dockerfile
LABEL author="Omar rady & Hazem elmahdy" version="1.0" description="ASCII Art web"
```

Read them with:

```bash
docker image inspect --format '{{ json .Config.Labels }}' ascii-art-web
```

### Unused objects

Stopped containers, dangling images and build cache accumulate over time. Clean them with:

```bash
docker rm ascii-art-web        # remove the stopped container
docker image prune -f          # remove dangling images
docker builder prune -f        # remove unused build cache
```

---

## Styling

The interface is a single dark-themed page styled in `static/style.css`.

- **Global reset** — a `*` rule removes the default margin and padding and applies `box-sizing: border-box`, so sizes behave the same across browsers
- **CSS custom properties** — colours and the font family are defined once in `:root` (`--bg-color`, `--surface-color`, `--text-color`, `--primary-color`, `--primary-hover`, `--font-family`) and reused everywhere, so the theme changes in one place
- **Responsive design** — a `@media (max-width: 600px)` rule reduces spacing and font sizes, and lets the banner options wrap on small screens
- **Interaction feedback** — the text area shows a green outline on focus, banner labels dim on hover, and the button changes colour smoothly on hover
- **Overflow** — the result box scrolls horizontally, so wide ASCII art never breaks the layout
- **HTML basics** — `lang="en"` on the root element and a viewport meta tag for mobile browsers

---

## Example

1. Open `http://localhost:8080` in a browser
2. Type `Hello` in the text field
3. Select the `shadow` banner
4. Click **Generate**
5. The server renders the ASCII art below the form

---

## Notes

- Only the Go standard library is used — no third-party packages
- The server listens on port `8080`; when using Docker, map it with `-p <host-port>:8080`
- Banner files, `templates/` and `static/` must stay in the same directory as the running binary — the `Dockerfile` copies them next to it
- The project is designed for local use; no authentication or HTTPS is included
