FROM golang:1.21-alpine AS build
LABEL authors="mohammadhoseinaref"

WORKDIR /app
COPY . .

RUN go build -o app ./


FROM alpine:3.18

WORKDIR /root/
COPY --from=build /app/app .

EXPOSE 8000
CMD ["./app"]
