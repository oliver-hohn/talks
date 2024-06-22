# Screen scraping in Go

## Setup
1. Install Go: https://go.dev/dl/ (skip if only running via Docker)
1. Install Google Chrome: https://www.google.com/intl/en_uk/chrome/dr/download/ (skip if only running via Docker)
1. Install Docker Desktop: https://docs.docker.com/desktop/

## Run locally
1. Run the script:
   ```
   go run main.go
   ```
   _`--headless` can be provided as an option to run in "headless" mode (without a UI)_

## Run via Docker
1. Build the container:
   ```
   docker build -t "screen_scraping_in_go:v1" .
   ```
1. Run the container:
   ```
   docker run screen_scraping_in_go:v1
   ```
   _runs in "headless" mode by default_
