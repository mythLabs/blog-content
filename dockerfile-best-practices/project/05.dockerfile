From node:18.17.1-alpine

RUN apk add --no-cache tini

RUN addgroup -S user && adduser -S user -G user

WORKDIR /app

COPY --chown=user:user . /app

COPY package.json package-lock.json .

RUN NODE_ENV=production && npm ci --production && npm cache clean --force

COPY . .

COPY data.json  .

USER user

EXPOSE 3000

CMD ["/sbin/tini", "--","node", "app.js" ]
