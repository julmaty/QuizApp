FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy module files and download modules during build (no vendor required)
COPY go.mod go.sum ./

# Copy the rest of the source
COPY . .

ENV CGO_ENABLED=0
# Ensure modules are downloaded inside the builder image
RUN go mod download
RUN go build -o /quizapp ./...

FROM alpine:3.18
COPY --from=builder /quizapp /quizapp
EXPOSE 8080
CMD ["/quizapp"]
