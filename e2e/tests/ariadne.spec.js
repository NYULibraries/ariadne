import * as fs from 'node:fs';

import { removeSourceMappingUrlComments, updateGoldenFiles } from '../testutils';

import { execSync } from 'child_process';
import { getTestCasesBackendSuccess } from '../../frontend/src/testutils';

const { test, expect } = require('@playwright/test');
const beautifyHtml = require('js-beautify').html;

const testCasesBackendSuccess = getTestCasesBackendSuccess();

for (let i = 0; i < testCasesBackendSuccess.length; i++) {
  const testCase = testCasesBackendSuccess[i];

  const stubBackendAPIResponse = async (page) => {
    // Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
    await page.route('**/v0/**', async (route) => {
      // Return a mock response with a JSON body and a 200 status code
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(testCase.response, null, '    '),
      });
    });
  };

  test.describe(`${testCase.name}`, () => {
    test.beforeEach(async ({ page }) => {
      await stubBackendAPIResponse(page);
      //await page.goto(`/${testCase.queryString}`);
      //await page.goto(`/`);
      await page.goto(`/?${testCase.queryString}`);
      //await history.pushState({}, null, `?${testCase.queryString}`);
      //await page.evaluate((queryString) => history.pushState({}, null, `?${queryString}`), testCase.queryString);
    });

    test('page HTML matches expected', async ({ page }) => {
      // Clean actual/ and diffs/ files
      // NOTE:
      // We don't bother with error handling because these files get overwritten
      // anyway, and if there were no previous files, or if a previous cleaning/reset
      // script or process already deleted the previous files, we don't want the errors
      // causing distraction.
      // If deletion fails on existing files, there's a good chance there will
      // be errors thrown later, which will then correctly fail the test.
      const actualFile = `tests/actual/${testCase.key}.html`;
      try {
        fs.unlinkSync(actualFile);
      } catch (error) { }
      const diffFile = `tests/diffs/${testCase.key}.txt`;
      try {
        fs.unlinkSync(diffFile);
      } catch (error) { }

      await page.waitForSelector('h6');

      const actual = beautifyHtml(removeSourceMappingUrlComments(await page.content()));

      const goldenFile = `tests/golden/${testCase.key}.html`;
      if (updateGoldenFiles()) {
        fs.writeFileSync(goldenFile, actual);

        console.log(`Updated golden file ${goldenFile}`);

        return;
      }
      const golden = beautifyHtml(fs.readFileSync(goldenFile, { encoding: 'utf8' }));

      fs.writeFileSync(actualFile, actual);

      const ok = actual === golden;

      let message = `Actual HTML for "${testCase.name}" does not match expected HTML`;
      if (!ok) {
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

    //test('screenshot matches expected', async ({ page }) => {
    //  // Wait for the response to be returned and the page to render
    //  await page.waitForSelector('.image');
    //  await page.waitForSelector('h6');

    //  // Take a screenshot to verify that the page was rendered correctly
    //  await expect(page).toHaveScreenshot(`${testCase.key}.png`);
    //});

    test('renders links with a "list-group" class name', async ({ page }) => {
      expect(await page.$('.list-group')).toBeTruthy();
    });

    test('returns search results', async ({ page }) => {
      expect(await page.textContent('p')).toBe('Displaying search results...');
    });
  });
}
