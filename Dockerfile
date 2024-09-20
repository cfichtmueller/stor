FROM node:22.9.0-alpine3.19 AS node
WORKDIR /stor
COPY . .
RUN npx tailwindcss -i ./internal/ui/css/input.css -o ./internal/ui/css/style.css -m

FROM golang:1.23.1-alpine3.19 AS build

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN apk add --no-cache gcc musl-dev
WORKDIR /stor
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY --from=node /stor/internal/ui/css/style.css ./internal/ui/css/style.css
COPY . .
RUN go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension -o stor main.go
RUN go test ./...

FROM scratch

VOLUME /var/data
EXPOSE 8000/tcp
EXPOSE 8001/tcp

COPY --from=build /stor/stor /usr/bin/stor

ENTRYPOINT [ "/usr/bin/stor" ]
CMD [ "serve" ]