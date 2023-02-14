const { test, expect } = require('@playwright/test');
const fs = require('fs');

import { getTestCasesBackendSuccess } from '../../frontend/src/testutils';

const testCasesBackendSuccess = getTestCasesBackendSuccess();

for (let i = 0; i < testCasesBackendSuccess.length; i++) {
  const testCase = testCasesBackendSuccess[i];
  const stubbedBackendAPIResponseBody = JSON.stringify(testCase.response, null, '    ');

  test.describe(`${testCase.name}`, () => {
    test.beforeEach(async() => {
      //Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
      await page.route('**/v0/*', async (route) => {
        //Return a mock response with a JSON body and a 200 status code
        await route.fulfill({
                              status: 200,
                              contentType: 'application/json',
                              body: stubbedBackendAPIResponseBody,
                            });
      });
    });

    test('HTML matches expected', async ({ page }) => {
      await page.goto(`/${testCase.queryString}`);
      await page.waitForSelector('h6');

      const snapshot = await page.innerHTML('body');
      // fs.writeFileSync('tests/actual/new-yorker.html', snapshot);
      const golden = fs.readFileSync(`tests/golden/${testCase.key}.html`, { encoding: 'utf8' });
      const ok = snapshot === golden;

      // Actual diffing is not working in CI:

      // if (!ok) {
      //   try {
      //     // eslint-disable-next-line no-unused-vars
      //     execSync('diff -c tests/golden/new-yorker.html tests/actual/new-yorker.html');
      //   } catch (error) {
      // TODO: save stdout and stderr into variables, and print them out in the test assertions
      //     // eslint-disable-next-line no-console
      //     console.log(error.stdout.toString());
      //     // eslint-disable-next-line no-console
      //     console.log(error.stderr.toString());
      //   }
      // }
      expect(ok).toBe(true);
    });

    test('Ask a Librarian link pops up a new Ask a Library tab', async ({ page }) => {
      await page.goto(`/${testCase.queryString}`);
      // Playwright's team recommendation for handling popups: https://playwright.dev/docs/pages#handling-popups
      // Start waiting for popup before clicking. Note no await.
      const popupPromise = page.waitForEvent('popup');
      await page.getByRole('link', { name: 'Ask a Librarian' }).click();
      const popup = await popupPromise;
      // Wait for the popup to load.
      await popup.waitForLoadState();

      expect(await page.textContent('.ask-librarian')).toBe('Need help?Ask a Librarian');
      expect(popup.url()).toBe('https://library.nyu.edu/ask/');
    });

    test('matches screenshot', async ({ page }) => {
      //Navigate to the frontend URL, which will make a request to the backend URL
      await page.goto(`/${testCase.queryString}`);

      //Wait for the response to be returned and the page to render
      await page.waitForSelector('.image');
      await page.waitForSelector('h6');

      //Take a screenshot to verify that the page was rendered correctly
      await expect(page).toHaveScreenshot(`${testCase.key}.png`);
    });

    test('renders links with a list-group className', async ({ page }) => {
      await page.goto(`/${testCase.queryString}`);
      expect(await page.$('.list-group')).toBeTruthy();
    });

    test('renders the search results', async ({ page }) => {
      await page.goto(`/${testCase.queryString}`);
      expect(await page.textContent('p')).toBe('Displaying search results...');
    });
  });
}
