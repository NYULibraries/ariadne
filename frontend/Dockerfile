FROM node:18.9.0-alpine

ARG production

ENV INSTALL_PATH /app

COPY frontend/package.json frontend/yarn.lock /tmp/
RUN cd /tmp && yarn install $(if [[ ! -z $production ]]; then echo "--production"; fi) \
  && mkdir -p ${INSTALL_PATH} \
  && cd ${INSTALL_PATH} \
  && cp -R /tmp/node_modules ${INSTALL_PATH} \
  && rm -r /tmp/* && yarn cache clean

WORKDIR ${INSTALL_PATH}

COPY frontend .

RUN chown -R node:node ${INSTALL_PATH}

USER node

RUN REACT_APP_API_URL=/v0/ yarn run build

EXPOSE 3000

CMD ["yarn", "start"]
