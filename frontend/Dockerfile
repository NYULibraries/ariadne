FROM node:18.9.0-alpine

ARG production

ENV INSTALL_PATH /app

COPY package.json yarn.lock /tmp/
RUN cd /tmp && yarn install $(if [[ ! -z $production ]]; then echo "--production"; fi) \
  && mkdir -p ${INSTALL_PATH} \
  && cd ${INSTALL_PATH} \
  && cp -R /tmp/node_modules ${INSTALL_PATH} \
  && rm -r /tmp/* && yarn cache clean

WORKDIR ${INSTALL_PATH}

COPY . .

EXPOSE 3000

CMD ["yarn", "start"]
