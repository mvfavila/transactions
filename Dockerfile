FROM golang:alpine AS build

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev linux-headers

COPY . /app

RUN go build -o transactions .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/transactions /app

COPY config/ ./config/
EXPOSE 8080

CMD ["./transactions"]
