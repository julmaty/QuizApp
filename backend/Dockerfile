FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy and build. We avoid installing git in the image to reduce network/package issues.
COPY go.mod .
COPY . .

ENV CGO_ENABLED=0
RUN go build -o /quizapp

FROM alpine:3.18
COPY --from=builder /quizapp /quizapp
EXPOSE 8080
CMD ["/quizapp"]
