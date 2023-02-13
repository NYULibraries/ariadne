const { test, expect } = require('@playwright/test');
const fs = require('fs');
const path = require('path');

const getQueryStrings = require('./utils/getQueryStrings');

const queryStrings = getQueryStrings();

const NYJSONDATA = fs.readFileSync(path.join(__dirname, '../../backend/api/testdata/server/golden/the-new-yorker.json'), { encoding: 'utf8' });
const CFJSONDATA = fs.readFileSync(path.join(__dirname, '../../backend/api/testdata/server/golden/corriere-fiorentino.json'), { encoding: 'utf8' });


test('Ask a Librarian link pops up a new Ask a Library tab', async ({ page }) => {
  await page.route('**/v0/*', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  await page.goto('/' + queryStrings[0]);
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


test('renders the New Yorker page', async ({ page }) => {
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
  await page.waitForSelector('.image');
  await page.waitForSelector('h6');

  //Take a screenshot to verify that the page was rendered correctly
  await expect(page).toHaveScreenshot('new_yorker.png');
});

test('renders the Corriere Fiorentino page', async ({ page }) => {
  // Define a mock HTTP request handler to intercept the request and log the details
  await page.route('**/v0/*', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: CFJSONDATA
    });
  });

  await page.goto('/' + queryStrings[1]);

  await page.waitForSelector('.image');
  await page.waitForSelector('h6', { timeout: 10000 });

  await expect(page).toHaveScreenshot('corriere_fiorentino.png');
});



test('renders links with a list-group className', async ({ page }) => {
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

test('renders the search results', async ({ page }) => {
  await page.route('**/v0/*', route => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: NYJSONDATA
    });
  });

  await page.goto('/' + queryStrings[0]);
  expect(await page.textContent('p')).toBe('Displaying search results...');
});


test('renders a E Journal Full Text link', async ({ page }) => {
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
  await page.waitForSelector('.image');
  await page.waitForSelector('h6');

  // Check if the expected URL is present in the DOM
  const link = await page.$('a[href="' + expectedUrl + '"]');
  expect(link).not.toBeNull();
});

test('renders the correct NYUAD logo and link based on institution query parameter', async ({ page }) => {
  await page.goto('/' + '?institution=NYUAD');
  const linkElement = await page.waitForSelector('a');
  expect(await linkElement.getAttribute('href')).toBe('https://nyuad.nyu.edu/en/library.html');
  const imgElement = await linkElement.$('img');
  expect(await imgElement.getAttribute('src')).toBe(`/images/abudhabi-logo-color.svg`);
});

test('renders the correct NYUSH logo and link based on institution query parameter', async ({ page }) => {
  await page.goto('/' + '?institution=NYUSH');
  const linkElement = await page.waitForSelector('a');
  expect(await linkElement.getAttribute('href')).toBe('https://shanghai.nyu.edu/academics/library');
  const imgElement = await linkElement.$('img');
  expect(await imgElement.getAttribute('src')).toBe(`/images/shanghai-logo-color.svg`);
});

test('redirects correctly when institution query parameter is "umlaut.institution"', async ({ page }) => {
  await page.goto('/' + '?umlaut.institution=NYSH');
  setTimeout(async () => {
    const linkElement = await page.waitForSelector('a', { timeout: 10000 });
    expect(await linkElement.getAttribute('href')).toBe('https://shanghai.nyu.edu/academics/library');
    const imgElement = await linkElement.$('img');
    expect(await imgElement.getAttribute('src')).toBe(`/images/shanghai-logo-color.svg`);
  }, 1000);
});


