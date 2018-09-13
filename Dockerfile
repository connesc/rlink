FROM golang:1.11 AS build
WORKDIR /app
COPY . .
RUN go build

FROM gcr.io/distroless/base
COPY --from=build /app/rlink /usr/local/bin/rlink
ENTRYPOINT [ "rlink" ]
