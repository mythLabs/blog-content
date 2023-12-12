From node:18.17.1-alpine

copy . .

RUN npm install

CMD [ "node", "app.js" ]
