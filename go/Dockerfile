FROM golang:1.24.0-alpine3.20

# 作業ディレクトリを作成
WORKDIR /app
COPY . .

# SQLite3ドライバのビルドにはCGOが必要
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev sqlite

RUN addgroup -S mercari && adduser -S trainee -G mercari
RUN chown -R trainee:mercari db
RUN chmod -R 775 db

# 依存関係のキャッシュを有効にするため、まず go.mod と go.sum をコピー
COPY go.mod go.sum ./
RUN go mod download  

USER trainee

CMD ["go", "run", "cmd/api/main.go"]