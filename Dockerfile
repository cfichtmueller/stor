FROM node:24.12.0-alpine3.23 AS node
WORKDIR /stor
COPY . .
RUN npx tailwindcss@3.4.19 -i ./internal/ui/css/input.css -o ./internal/ui/css/style.css -m

FROM golang:1.25.5-alpine3.23 AS build

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