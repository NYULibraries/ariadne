FROM mcr.microsoft.com/playwright:v1.27.1-focal

WORKDIR /e2e

COPY . /e2e

# Install dependencies
RUN yarn install

RUN npx playwright install

# RUN playwright tests
CMD ["yarn", "test:e2e"]



