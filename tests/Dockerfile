FROM node:18.16.0-alpine3.17

WORKDIR /srv

COPY package* ./

RUN npm ci

COPY . ./

CMD ["npm", "run", "test"]