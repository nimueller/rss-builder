FROM golang:1.25-alpine AS build
WORKDIR /build
COPY *.go go.mod go.sum /build/
RUN CGO_ENABLED=0 go build .

FROM alpine:3.23.2
WORKDIR /app
COPY --from=build /build/rss-builder /app/
ENV DATABASE_URL=""
ENV WEBSERVER_HOST=0.0.0.0
ENV WEBSERVER_PORT=8080
ENV SCRAPER_INTERVAL=10m
EXPOSE 8080
CMD [ "/app/rss-builder" ]

