const { test, expect } = require('@playwright/test');
const fs = require('fs');
const path = require('path');

import {
  getTestCasesBackendSuccess,
} from '../../frontend/src/testutils';

const testCasesBackendSuccess = getTestCasesBackendSuccess();

for (let i = 0; i < testCasesBackendSuccess.length; i++) {
  const testCase = testCasesBackendSuccess[ i ];
  const stubbedBackendAPIResponse = JSON.stringify(testCase.response, null, '    ');

  test.describe(`${testCase.name}`, () => {

    test('Ask a Librarian link pops up a new Ask a Library tab', async ({ page }) => {
      await page.route('**/v0/*', async route => {
        await route.fulfill({
                              status: 200,
                              contentType: 'application/json',
                              body: stubbedBackendAPIResponse
                            });
      });

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
      //Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
      await page.route('**/v0/*', async route => {
        //Return a mock response with a JSON body and a 200 status code
        await route.fulfill({
                              status: 200,
                              contentType: 'application/json',
                              body: stubbedBackendAPIResponse
                            });
      });

      //Navigate to the frontend URL, which will make a request to the backend URL
      await page.goto(`/${testCase.queryString}`);

      //Wait for the response to be returned and the page to render
      await page.waitForSelector('.image');
      await page.waitForSelector('h6');

      //Take a screenshot to verify that the page was rendered correctly
      await expect(page).toHaveScreenshot(`${testCase.key}.png`);
    });

    test('renders links with a list-group className', async ({ page }) => {
      await page.route('**/v0/*', route => {
                         route.fulfill({
                                         status: 200,
                                         contentType: 'application/json',
                                         body: stubbedBackendAPIResponse
                                       });
                       }
      );

      await page.goto(`/${testCase.queryString}`);
      expect(await page.$('.list-group')).toBeTruthy();
    });

    test('renders the search results', async ({ page }) => {
      await page.route('**/v0/*', route => {
        route.fulfill({
                        status: 200,
                        contentType: 'application/json',
                        body: stubbedBackendAPIResponse
                      });
      });

      await page.goto(`/${testCase.queryString}`);
      expect(await page.textContent('p')).toBe('Displaying search results...');
    });

  });
}

