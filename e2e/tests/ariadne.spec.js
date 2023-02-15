import { execSync } from 'child_process';
const fs = require('fs');

const { test, expect } = require('@playwright/test');
const beautifyHtml = require('js-beautify').html;

import { getTestCasesBackendSuccess } from '../../frontend/src/testutils';
import { updateGoldenFiles } from '../testutils';

const testCasesBackendSuccess = getTestCasesBackendSuccess();

for (let i = 0; i < testCasesBackendSuccess.length; i++) {
  const testCase = testCasesBackendSuccess[i];

  const stubBackendAPIResponse = async (page) => {
    //Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
    await page.route('**/v0/*', async (route) => {
      //Return a mock response with a JSON body and a 200 status code
      await route.fulfill({
                            status: 200,
                            contentType: 'application/json',
                            body: JSON.stringify(testCase.response, null, '    '),
                          });
    });
  }

  test.describe(`${testCase.name}`, () => {
    test.beforeEach(async ({ page }) => {
      await stubBackendAPIResponse(page);
      await page.goto(`/${testCase.queryString}`);
    });

    test('HTML matches expected', async ({ page }) => {
      await page.waitForSelector('h6');

      const actualFile = `tests/actual/${testCase.key}.html`;
      const actual = beautifyHtml(await page.innerHTML('body'));
      const goldenFile = `tests/golden/${testCase.key}.html`;
      const golden = beautifyHtml(fs.readFileSync(goldenFile, { encoding: 'utf8' }));

      if ( updateGoldenFiles() ) {
        fs.writeFileSync( goldenFile, actual );

        console.log( `Updated golden file ${ goldenFile }` );

        return;
      }

      fs.writeFileSync(actualFile, actual);

      const ok = ( actual === golden );

      let message = 'Actual HTML for `${testCase.name}` does not match expected.';
      if (!ok) {
        const diffFile = `tests/diffs/${testCase.key}.txt`;
        const command = `diff ${goldenFile} ${actualFile} | tee ${diffFile}`;
        let diffOutput;
        try {
          diffOutput = new TextDecoder().decode(execSync(command));
          message += `

======= BEGIN DIFF OUTPUT ========
${diffOutput}
======== END DIFF OUTPUT =========

[Recorded in diff file: ${diffFile}]`;
        } catch (e) {
          // `diff` command failed to create the diff file.
          message += `  Diff command \`${command}\` failed:

${e.stderr.toString()}`;
        }
      }

      expect(ok, message).toBe(true);
    });

    test('Ask a Librarian link pops up a new Ask a Library tab', async ({ page }) => {
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
      //Wait for the response to be returned and the page to render
      await page.waitForSelector('.image');
      await page.waitForSelector('h6');

      //Take a screenshot to verify that the page was rendered correctly
      await expect(page).toHaveScreenshot(`${testCase.key}.png`);
    });

    test('renders links with a list-group className', async ({ page }) => {
      expect(await page.$('.list-group')).toBeTruthy();
    });

    test('renders the search results', async ({ page }) => {
      expect(await page.textContent('p')).toBe('Displaying search results...');
    });
  });
}
