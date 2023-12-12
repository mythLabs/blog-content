From node:latest

copy . .

RUN npm install

CMD [ "node", "app.js" ]
