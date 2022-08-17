FROM golang:1.19-alpine3.16 AS dev

RUN apk add --no-cache \
    build-base \
    gcc \
    git
    
RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

COPY ./ /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build main.go
ENTRYPOINT /go/bin/CompileDaemon --build="go build main.go"

FROM alpine:3.16 AS prod
WORKDIR /app
COPY --from=dev /app/main .
CMD [ "./main" ]