FROM node:20-alpine

WORKDIR /app
COPY package*.json ./

RUN npm ci

COPY . .


RUN addgroup -S mercari && adduser -S trainee -G mercari

# node_modulesディレクトリの権限を調整
RUN chown -R trainee:mercari /app/node_modules && \
    chmod -R 755 /app/node_modules

# コンテナ内で外部アクセスするために使用するポートを明示
EXPOSE 3000

USER trainee

# Vite の起動時に --host 0.0.0.0 を付与して、外部アクセスを許可
CMD ["npm", "start", "--", "--host", "0.0.0.0"]