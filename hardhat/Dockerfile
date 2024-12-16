FROM node:20-alpine

WORKDIR /usr/src/app

RUN apk add --no-cache bash

COPY package.json  ./

RUN yarn install --frozen-lockfile

COPY . .

EXPOSE 8545

CMD ["yarn", "start-node"]