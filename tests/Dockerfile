FROM node:17.6.0-alpine3.15

WORKDIR /srv

COPY package* ./

RUN npm ci

COPY . ./

CMD ["npm", "run", "test"]