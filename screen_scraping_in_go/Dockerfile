# ref: https://docs.docker.com/language/golang/build-images/#multi-stage-builds
FROM golang:1.22 AS build_stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/scraper

FROM alpine:3.19 AS release_stage
WORKDIR /app

# Install Chrome. From: https://github.com/Zenika/alpine-chrome/blob/ab0d418bcbc81e1254bdb7875ec06753e3fd68c5/Dockerfile
# License: Apache License 2.0
# Changed to only install Chrome, and not run it by default when running the container.
RUN apk upgrade --no-cache --available \
  && apk add --no-cache \
  chromium-swiftshader \
  ttf-freefont \
  font-noto-emoji \
  && apk add --no-cache \
  --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community \
  font-wqy-zenhei

ENV CHROME_BIN=/usr/bin/chromium-browser \
  CHROME_PATH=/usr/lib/chromium/ \
  # Autorun chrome headless
  CHROMIUM_FLAGS="--disable-software-rasterizer --disable-dev-shm-usage"

# Copy all built binaries across
COPY --from=build_stage /tmp/scraper /app/scraper

# Create a system user to not run any of the binaries as root
RUN addgroup -g 10001 -S appgroup && \
  adduser -u 10000 -S -G appgroup appuser && \
  chown -R appuser:appgroup /app/

USER appuser
ENTRYPOINT [ "/app/scraper", "--headless" ]
