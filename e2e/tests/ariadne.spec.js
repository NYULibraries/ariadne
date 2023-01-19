const { test, expect } = require('@playwright/test');


const queryStrings = [
  '?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat',
  '?ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=%3CNYUMARCIT%3E3400000000000901%3C/NYUMARCIT%3E%3Cgrp_id%3E582323038%3C/grp_id%3E%3Coa%3E%3C/oa%3E%3Curl%3E%3C/url%3E&rft_id=info:oai/&req.language=eng',
];

test('renders the Corriere Fiorentino page correctly', async ({ page }) => {
  await page.goto( '/' + queryStrings[1]);

  await page.waitForFunction(() => document.querySelector('.image'));
  await page.waitForFunction(() => document.querySelector('h6'), { timeout: 10000 });

  await expect(page).toHaveScreenshot('corriere_fiorentino.png');
});

test.skip('renders a Press Reader link', async ({ page }) => {
  await page.goto('/' + queryStrings[1]);
  await page.waitForFunction(() => document.querySelector('h6'));

  const [page2] = await Promise.all([
    page.waitForEvent('popup'),
    page.getByRole('link', { name: 'PressReader' }).click(),
  ]);
  expect(page2.url()).toBe('https://www.pressreader.com/italy/corriere-fiorentino');
});

test('renders the New Yorker page correctly', async ({ page }) => {
  await page.goto('/' + queryStrings[0]);

  await page.waitForFunction(() => document.querySelector('.image'));
  await page.waitForFunction(() => document.querySelector('h6'));

  await expect(page).toHaveScreenshot('new_yorker.png');
});

test('renders with a className of list-group', async ({ page }) => {
  await page.goto('/'+ queryStrings[0]);
  expect(await page.$('.list-group')).toBeTruthy();
});

test('renders the search results', async ({ page }) => {
  await page.goto('/'+ queryStrings[0]);
  expect(await page.textContent('p')).toBe('Displaying search results...');
});

test('renders Loading...', async ({ page }) => {
  await page.goto('/' + queryStrings[0]);
  expect(await page.textContent('.loader')).toBe('Loading...');
});

test.skip('renders a E Journal Full Text link', async ({ page }) => {
  await page.goto('/'+ queryStrings[0]);
  const [page2] = await Promise.all([
    page.waitForEvent('popup'),
    page.getByRole('link', { name: 'E Journal Full Text' }).click(),
  ]);
  expect(page2.url()).toBe('https://archives.newyorker.com/#folio=C1');
});

test.skip('renders a Gale General OneFile link', async ({ page }) => {
  await page.goto('/' + queryStrings[0]);
  const [page3] = await Promise.all([
    page.waitForEvent('popup'),
    page.getByRole('link', { name: 'Gale General OneFile' }).click(),
  ]);
  expect(page3.url()).toBe(
    'https://go.gale.com/ps/i.do?p=ITOF&u=nysl_me_newyorku&id=GALE%7C1161&v=2.1&it=aboutJournal'
  );
});

test('renders Ask a Librarian', async ({ page }) => {
  await page.goto('/'+ queryStrings[0]);
  const [page4] = await Promise.all([
    page.waitForEvent('popup', { timeout: 10000 }),
    page.getByRole('link', { name: 'Ask a Librarian' }).click(),
  ]);
  expect(await page.textContent('.ask-librarian')).toBe('Need help?Ask a Librarian');
  expect(page4.url()).toBe('https://library.nyu.edu/ask/');
});

test('Loading... no longer present in the DOM after loading data', async ({ page }) => {
  await page.goto('/' + queryStrings[0]);
  expect(await page.textContent('p')).not.toBe('Loading...');
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
    const linkElement = await page.waitForSelector('a');
    expect(await linkElement.getAttribute('href')).toBe('https://shanghai.nyu.edu/academics/library');
    const imgElement = await linkElement.$('img');
    expect(await imgElement.getAttribute('src')).toBe(`/images/shanghai-logo-color.svg`);
  }, 1000);
});
  

  