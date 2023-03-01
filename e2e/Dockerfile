FROM mcr.microsoft.com/playwright:v1.27.1-focal

WORKDIR /e2e

COPY e2e/package.json e2e/yarn.lock ./

# Install dependencies
RUN yarn install

RUN npx playwright install

# copy over remaining e2e files
COPY e2e/. .

COPY backend/api/testdata/server/golden/* /backend/api/testdata/server/golden/

COPY backend/api/testdata/server/test-cases.json /backend/api/testdata/server/test-cases.json

COPY frontend/src/testutils/index.js /frontend/src/testutils/index.js


# entrypoint ensures $PLAYWRIGHT_BASE_URL has started
ENTRYPOINT [ "./docker-entrypoint.sh" ]
# RUN playwright tests
CMD ["yarn", "test:e2e"]
