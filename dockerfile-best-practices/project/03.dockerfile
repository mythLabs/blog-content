From node:18.17.1-alpine

WORKDIR /app

COPY package.json package-lock.json .

RUN NODE_ENV=production && npm ci --production && npm cache clean --force

COPY . .

COPY data.json  .

CMD [ "node", "app.js" ]
