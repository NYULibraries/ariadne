FROM node:18.9.0-alpine

ARG production

RUN addgroup --gid 1001 docker && \
  adduser --uid 1001 --ingroup docker --disabled-password --gecos "" docker

ENV INSTALL_PATH /app

COPY frontend/package.json frontend/yarn.lock /tmp/
RUN cd /tmp && yarn install $(if [[ ! -z $production ]]; then echo "--production"; fi) \
  && mkdir -p ${INSTALL_PATH} \
  && cd ${INSTALL_PATH} \
  && cp -R /tmp/node_modules ${INSTALL_PATH} \
  && rm -r /tmp/* && yarn cache clean

WORKDIR ${INSTALL_PATH}

COPY frontend .

RUN REACT_APP_API_URL=/v0/ yarn run build

USER docker

EXPOSE 3000

CMD ["yarn", "start"]
