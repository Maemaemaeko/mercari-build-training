FROM node:20-alpine

WORKDIR /app
COPY package*.json ./

RUN npm ci

COPY . .


RUN addgroup -S mercari && adduser -S trainee -G mercari

# node_modulesディレクトリの権限を調整
RUN chown -R trainee:mercari /app/node_modules && \
    chmod -R 755 /app/node_modules


USER trainee

CMD ["npm", "start"]
