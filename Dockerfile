FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY app.html app.go downloader.go  ./

RUN go build -o main .

FROM alpine

WORKDIR /app

COPY --from=build /app/main ./

COPY --from=build /app/app.html ./ 

ENTRYPOINT ["./main"]

