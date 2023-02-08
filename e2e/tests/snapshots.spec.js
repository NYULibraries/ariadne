import { expect, test } from '@playwright/test';

import fs from 'fs';

// This import is not working in CI
// import { execSync } from 'child_process';
const getQueryStrings = require('./utils/getQueryStrings');

const queryStrings = getQueryStrings();

const NYJSONDATA = fs.readFileSync(require('path').join(__dirname, '../../backend/api/testdata/server/golden/the-new-yorker.json'), { encoding: 'utf8' });
const CFJSONDATA = fs.readFileSync(require('path').join(__dirname, '../../backend/api/testdata/server/golden/corriere-fiorentino.json'), { encoding: 'utf8' });


test('compares the rendered Corriere Fiorentino page to a golden file', async ({ page }) => {
  //Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
  await page.route('**/v0/*', async route => {
    //Return a mock response with a JSON body and a 200 status code
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: CFJSONDATA
    });
  });

  //Navigate to the frontend URL, which will make a request to the backend URL
  await page.goto('/' + queryStrings[1]);
  await page.waitForFunction(() => document.querySelector('h6'));

  const snapshot = await page.innerHTML('body');
  // fs.writeFileSync('tests/actual/corriere-fiorentino.html', snapshot);
  const golden = fs.readFileSync('tests/golden/corriere-fiorentino.html', 'utf8');
  const stringifiedSnapshot = JSON.stringify(snapshot);
  const stringifiedGolden = JSON.stringify(golden);
  const ok = stringifiedSnapshot === stringifiedGolden;

  // Actual diffing is not working in CI:

  // if (stringifiedGolden !== stringifiedSnapshot) {
  //   try {
  //     execSync('diff -c tests/golden/corriere-fiorentino.html e2e/actual/corriere-fiorentino.html');
  //   } catch (error) {
  //     // eslint-disable-next-line no-console
  //     console.log(error.stdout.toString());
  //     // eslint-disable-next-line no-console
  //     console.log(error.stderr.toString());
  //   }
  // }

  expect(ok).toBeTruthy();
});

test('compares the rendered New Yorker page to a golden file', async ({ page }) => {
  //Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
  await page.route('**/v0/*', async route => {
    //Return a mock response with a JSON body and a 200 status code
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  await page.goto('/' + queryStrings[0]);
  await page.waitForFunction(() => document.querySelector('h6'));

  const snapshot = await page.innerHTML('body');
  // fs.writeFileSync('tests/actual/new-yorker.html', snapshot);
  const golden = fs.readFileSync('tests/golden/new-yorker.html', { encoding: 'utf8' });
  const stringifiedSnapshot = JSON.stringify(snapshot);
  const stringifiedGolden = JSON.stringify(golden);
  const ok = stringifiedSnapshot === stringifiedGolden;

  // Acrual diffing is not working in CI:

  // if (stringifiedGolden !== stringifiedSnapshot) {
  //   try {
  //     // eslint-disable-next-line no-unused-vars
  //     execSync('diff -c tests/golden/new-yorker.html e2e/actual/new-yorker.html');
  //   } catch (error) {
  //     // eslint-disable-next-line no-console
  //     console.log(error.stdout.toString());
  //     // eslint-disable-next-line no-console
  //     console.log(error.stderr.toString());
  //   }
  // }
  expect(ok).toBe(true);
});
