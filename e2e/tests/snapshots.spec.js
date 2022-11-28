import { test, expect } from '@playwright/test';
import fs from 'fs';
import { execSync } from 'child_process';

const queryStrings = [
  '?url_ver=Z39.88-2004&url_ctx_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Actx&ctx_ver=Z39.88-2004&ctx_tim=2021-10-22T12%3A29%3A27-04%3A00&ctx_id=&ctx_enc=info%3Aofi%2Fenc%3AUTF-8&rft.aulast=Ross&rft.date=2002&rft.eissn=2163-3827&rft.genre=journal&rft.issn=0028-792X&rft.jtitle=New+Yorker&rft.language=eng&rft.lccn=++2011201780&rft.object_id=110975413975944&rft.oclcnum=909782404&rft.place=New+York&rft.private_data=909782404%3Cfssessid%3E0%3C%2Ffssessid%3E&rft.pub=F-R+Pub.+Corp.&rft.stitle=NEW+YORKER&rft.title=New+Yorker&rft_val_fmt=info%3Aofi%2Ffmt%3Akev%3Amtx%3Ajournal&rft_id=info%3Aoclcnum%2F909782404&rft_id=urn%3AISSN%3A0028-792X&req.ip=209.150.44.95&rfr_id=info%3Asid%2FFirstSearch%3AWorldCat',
  '?ctx_ver=Z39.88-2004&ctx_enc=info:ofi/enc:UTF-8&ctx_tim=2018-07-15T02:13:26IST&url_ver=Z39.88-2004&url_ctx_fmt=infofi/fmt:kev:mtx:ctx&rfr_id=info:sid/primo.exlibrisgroup.com:primo-dedupmrg524707060&rft_val_fmt=info:ofi/fmt:kev:mtx:journal&rft.genre=journal&rft.jtitle=Corriere%20Fiorentino&rft.btitle=Corriere%20Fiorentino&rft.aulast=&rft.aufirst=&rft.auinit=&rft.auinit1=&rft.auinitm=&rft.ausuffix=&rft.au=&rft.aucorp=&rft.volume=&rft.issue=&rft.part=&rft.quarter=&rft.ssn=&rft.spage=&rft.epage=&rft.pages=&rft.artnum=&rft.pub=&rft.place=Italy&rft.issn=&rft.eissn=&rft.isbn=&rft.sici=&rft.coden=&rft_id=info:doi/&rft.object_id=3400000000000901&rft.primo=dedupmrg524707060&rft.eisbn=&rft_dat=%3CNYUMARCIT%3E3400000000000901%3C/NYUMARCIT%3E%3Cgrp_id%3E582323038%3C/grp_id%3E%3Coa%3E%3C/oa%3E%3Curl%3E%3C/url%3E&rft_id=info:oai/&req.language=eng',
];

test('renders the Corriere Fiorentino against golden correctly', async ({ page }) => {
  await page.goto('/' + queryStrings[1]);
  await page.waitForFunction(() => document.querySelector('h6'));

  const snapshot = await page.innerHTML('body');
  // fs.writeFileSync('tests/actual/corriere-fiorentino.html', snapshot);
  const golden = fs.readFileSync('tests/golden/corriere-fiorentino.html', 'utf8');
  const stringifiedSnapshot = JSON.stringify(snapshot);
  const stringifiedGolden = JSON.stringify(golden);
  const ok = stringifiedSnapshot === stringifiedGolden;

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

test.skip('renders the New Yorker against golden correctly', async ({ page }) => {
  await page.goto('/'+ queryStrings[0]);
  await page.waitForFunction(() => document.querySelector('h6'));

  const snapshot = await page.innerHTML('body');
  fs.writeFileSync('e2e/test/new-yorker1.html', snapshot);
  const golden = fs.readFileSync('tests/golden/new-yorker.html', { encoding: 'utf8' });
  const stringifiedSnapshot = JSON.stringify(snapshot);
  const stringifiedGolden = JSON.stringify(golden);
  const ok = stringifiedSnapshot === stringifiedGolden;

  if (stringifiedGolden !== stringifiedSnapshot) {
    try {
      // eslint-disable-next-line no-unused-vars
      execSync('diff -c tests/golden/new-yorker.html e2e/actual/new-yorker.html');
    } catch (error) {
      // eslint-disable-next-line no-console
      console.log(error.stdout.toString());
      // eslint-disable-next-line no-console
      console.log(error.stderr.toString());
    }
  }
  expect(ok).toBe(true);
});
