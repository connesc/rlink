FROM golang:1.11-alpine3.8 AS build
RUN apk add --no-cache git
WORKDIR /usr/src/rlink
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-s -w' -o /rlink

FROM scratch
COPY --from=build /rlink /
ENTRYPOINT [ "/rlink" ]
