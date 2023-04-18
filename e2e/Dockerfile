FROM mcr.microsoft.com/playwright:v1.27.1-focal

WORKDIR /e2e

COPY e2e/package.json e2e/yarn.lock ./

# Install dependencies
RUN yarn install

RUN npx playwright install

# copy over remaining e2e files
COPY e2e/. .

COPY backend/testutils/testdata/golden/api-responses/* /backend/testutils/testdata/golden/api-responses/

COPY backend/testutils/testdata/test-cases.json /backend/testutils/testdata/test-cases.json

COPY frontend/src/testutils/index.js /frontend/src/testutils/index.js


# entrypoint ensures $PLAYWRIGHT_BASE_URL has started
ENTRYPOINT [ "./docker-entrypoint.sh" ]
# RUN playwright tests
CMD ["yarn", "test:e2e"]
