FROM node:11.1-alpine

COPY . /project/swagger-ui

WORKDIR /project/swagger-ui

RUN npm install

EXPOSE 8080

CMD ["node", "index.js"]