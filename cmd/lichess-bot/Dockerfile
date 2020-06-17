FROM golang:1.14 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && \
    go get \
    golang.org/x/tools/cmd/stringer \
    github.com/matryer/moq

COPY . . 
RUN go generate ./...
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o lichess-bot ./cmd/lichess-bot/ 

FROM gcr.io/distroless/static

COPY --from=builder /build/lichess-bot .

ENTRYPOINT [ "/lichess-bot" ]
