const { test, expect } = require('@playwright/test');
const fs = require('fs');
const path = require('path');

const getQueryStrings = require('./utils/getQueryStrings');

const queryStrings = getQueryStrings();

const NYJSONDATA = fs.readFileSync(path.join(__dirname, '../../backend/api/testdata/server/golden/the-new-yorker.json'), { encoding: 'utf8'});
const CFJSONDATA = fs.readFileSync(path.join(__dirname, '../../backend/api/testdata/server/golden/corriere-fiorentino.json'), { encoding: 'utf8'});


test('stubbing out Ask a Librarian', async ({ page }) => {
  await page.route('**/v0/*', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  await page.goto('/' + queryStrings[0]);
  const [page4] = await Promise.all([
    page.waitForEvent('popup', { timeout: 10000 }),
    page.getByRole('link', { name: 'Ask a Librarian' }).click(),
  ]);
  expect(await page.textContent('.ask-librarian')).toBe('Need help?Ask a Librarian');
  expect(page4.url()).toBe('https://library.nyu.edu/ask/');
});


test('stubs out the New Yorker page request', async ({ page }) => {
  //Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
  await page.route('**/v0/*', async route => {
    //Return a mock response with a JSON body and a 200 status code
   await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  //Navigate to the frontend URL, which will make a request to the backend URL
  await page.goto('/' + queryStrings[0]);

  //Wait for the response to be returned and the page to render
  await page.waitForFunction(() => document.querySelector('.image'));
  await page.waitForFunction(() => document.querySelector('h6'));

  //Take a screenshot to verify that the page was rendered correctly
  await expect(page).toHaveScreenshot('new_yorker.png');
});

test('stubs out the Corriere Fiorentino page request', async ({ page }) => {
  // Define a mock HTTP request handler to intercept the request and log the details
  await page.route('**/v0/*', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: CFJSONDATA
    });
  });

  // Log messages from the browser console to the test console for debugging:
  // page.on('console', msg => console.log(msg.text()));

  await page.goto('/' + queryStrings[1]);

  await page.waitForFunction(() => document.querySelector('.image'));
  await page.waitForFunction(() => document.querySelector('h6'), { timeout: 10000 });

  await expect(page).toHaveScreenshot('corriere_fiorentino.png');
});



test('stubs out with a className of list-group', async ({ page }) => {
  await page.route('**/v0/*', route => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  }
  );

  await page.goto('/' + queryStrings[0]);
  expect(await page.$('.list-group')).toBeTruthy();
});

test('stubs out the search results', async ({ page }) => {
  await page.route('**/v0/*', route => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  await
  page.goto('/' + queryStrings[0]);
  expect(await page.textContent('p')).toBe('Displaying search results...');
});


test('stubs out a E Journal Full Text link', async ({ page }) => {
  // Replace this with the expected URL from the JSON file
  const expectedUrl = 'http://proxy.library.nyu.edu/login?url=http://archives.newyorker.com/#folio=C1';

  // Define a mock HTTP request handler for the /v0/ URL path to intercept the request and return a mocked response.
  await page.route('**/v0/*', async route => {
    // Return a mock response with a JSON body and a 200 status code
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  // Navigate to the frontend URL, which will make a request to the backend URL
  await page.goto('/' + queryStrings[0]);

  // Wait for the response to be returned and the page to render
  await page.waitForFunction(() => document.querySelector('.image'));
  await page.waitForFunction(() => document.querySelector('h6'));

  // Check if the expected URL is present in the DOM
  const link = await page.$('a[href="' + expectedUrl + '"]');
  expect(link).not.toBeNull();
});

// Test skipped because they use the actual backend API:

// test.skip('renders the Corriere Fiorentino page correctly', async ({ page }) => {
//   await page.goto( '/' + queryStrings[1]);

//   await page.waitForFunction(() => document.querySelector('.image'));
//   await page.waitForFunction(() => document.querySelector('h6'), { timeout: 10000 });

//   await expect(page).toHaveScreenshot('corriere_fiorentino.png');
// });

// test.skip('renders a Press Reader link', async ({ page }) => {
//   await page.goto('/' + queryStrings[1]);
//   await page.waitForFunction(() => document.querySelector('h6'));

//   const [page2] = await Promise.all([
//     page.waitForEvent('popup'),
//     page.getByRole('link', { name: 'PressReader' }).click(),
//   ]);
//   expect(page2.url()).toBe('https://www.pressreader.com/italy/corriere-fiorentino');
// });

// test.skip('renders the New Yorker page correctly', async ({ page }) => {
//   await page.goto('/' + queryStrings[0]);

//   await page.waitForFunction(() => document.querySelector('.image'));
//   await page.waitForFunction(() => document.querySelector('h6'));

//   await expect(page).toHaveScreenshot('new_yorker.png');
// });

// test.skip('renders with a className of list-group', async ({ page }) => {
//   await page.goto('/'+ queryStrings[0]);
//   expect(await page.$('.list-group')).toBeTruthy();
// });

// test.skip('renders the search results', async ({ page }) => {
//   await page.goto('/'+ queryStrings[0]);
//   expect(await page.textContent('p')).toBe('Displaying search results...');
// });

// test.skip('renders Loading...', async ({ page }) => {
//   await page.goto('/' + queryStrings[0]);
//   expect(await page.textContent('.loader')).toBe('Loading...');
// });

// test.skip('renders a E Journal Full Text link', async ({ page }) => {
//   await page.goto('/'+ queryStrings[0]);
//   const [page2] = await Promise.all([
//     page.waitForEvent('popup'),
//     page.getByRole('link', { name: 'E Journal Full Text' }).click(),
//   ]);
//   expect(page2.url()).toBe('https://archives.newyorker.com/#folio=C1');
// });

// test.skip('renders a Gale General OneFile link', async ({ page }) => {
//   await page.goto('/' + queryStrings[0]);
//   const [page3] = await Promise.all([
//     page.waitForEvent('popup'),
//     page.getByRole('link', { name: 'Gale General OneFile' }).click(),
//   ]);
//   expect(page3.url()).toBe(
//     'https://go.gale.com/ps/i.do?p=ITOF&u=nysl_me_newyorku&id=GALE%7C1161&v=2.1&it=aboutJournal'
//   );
// });

// test.skip('renders Ask a Librarian', async ({ page }) => {
//   await page.goto('/'+ queryStrings[0]);
//   const [page4] = await Promise.all([
//     page.waitForEvent('popup', { timeout: 10000 }),
//     page.getByRole('link', { name: 'Ask a Librarian' }).click(),
//   ]);
//   expect(await page.textContent('.ask-librarian')).toBe('Need help?Ask a Librarian');
//   expect(page4.url()).toBe('https://library.nyu.edu/ask/');
// });

// test.skip('Loading... no longer present in the DOM after loading data', async ({ page }) => {
//   await page.goto('/' + queryStrings[0]);
//   expect(await page.textContent('p')).not.toBe('Loading...');
// });
