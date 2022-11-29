FROM golang:1.19-alpine as build

WORKDIR /ws
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /out/app .

FROM alpine:latest
WORKDIR /ws
COPY --from=build /out/app /ws/app

EXPOSE 3000
ENTRYPOINT ["/ws/app","serve"]