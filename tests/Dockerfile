FROM node:20.11.0-alpine3.19

WORKDIR /srv

COPY package* ./

RUN npm ci

COPY . ./

CMD ["npm", "run", "test"]