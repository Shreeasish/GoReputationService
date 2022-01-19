# Multi-stage layout
FROM golang:1.17 as builder

ENV GO111MODULE=on
ENV GOPROXY direct

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
ADD  resources/ /app/resources/
COPY --from=builder /app/reputation /app/
EXPOSE 8080
ENTRYPOINT ["/app/reputation"]
