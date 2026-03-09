FROM golang:1.26-alpine
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod ./
RUN go mod download || true
EXPOSE 8080
CMD ["air"]
