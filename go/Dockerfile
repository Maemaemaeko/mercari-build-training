FROM golang:1.24.0-alpine3.20

ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev sqlite

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# ユーザー・グループの作成と権限設定
RUN addgroup -S mercari && adduser -S trainee -G mercari \
    && chown -R trainee:mercari db images

USER root

CMD ["go", "run", "cmd/api/main.go"]
