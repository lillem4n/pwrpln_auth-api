FROM node:16.8.0-slim

WORKDIR /srv

COPY package* ./

RUN npm ci

COPY . ./

CMD ["npm", "run", "test"]