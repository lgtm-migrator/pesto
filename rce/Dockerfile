FROM node:16.14.0-bullseye

ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /app

RUN apt-get update && \
  apt-get upgrade -y && \
  apt-get install -y coreutils binutils build-essential libseccomp-dev gcc apt-utils curl tar bzip2 gzip

COPY ./nosocket ./nosocket

RUN make -C ./nosocket/ all && make -C ./nosocket/ install

COPY ./packages ./packages

COPY ./scripts ./scripts

COPY package*.json .

RUN npm ci

RUN node ./scripts/install.cjs

RUN node ./scripts/register-users.cjs

COPY tsconfig.json .

COPY ./src ./src

COPY build.js .

RUN npm run build

RUN rm -rf node_modules && npm ci --omit=dev

ENV PORT=50051

ENV NODE_ENV=production

EXPOSE ${PORT}

CMD ["node", "./dist/index.js"]
