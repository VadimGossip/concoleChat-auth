##
## Stage 1
##
FROM golang:1.22-alpine AS build

WORKDIR /build
ADD . /build
RUN go build -o app cmd/main.go

##
## Stage 2
##
FROM alpine:latest

WORKDIR /app
COPY --from=build /build/app /app/auth
COPY --from=build /build/config/config.yml /app/config/config.yml

EXPOSE $APP_PORT
CMD ["./auth"]