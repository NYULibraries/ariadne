import { getParameterFromQueryString } from './helpers';
import { bannerInstitutionInfo } from './institutionConstants';

const institutionNamesUpperCase = Object.keys(bannerInstitutionInfo).map(institutionName => institutionName.toUpperCase());
describe.each(institutionNamesUpperCase)(
    'Institution name: %s', (institutionName) => {
      test(`returns the correct value for "institution=${institutionName}" query parameter`, () => {
          const queryString = `?institution=${institutionName}`;
          const parameterName = 'institution';
          const expectedValue = institutionName.toLowerCase();
          const returnedValue = getParameterFromQueryString(queryString, parameterName);
          expect(returnedValue).toBe(expectedValue);
      });

      test(`returns the correct value for "institution=${institutionName.toLowerCase()}" query parameter`, () => {
        const queryString = `?institution=${institutionName.toLowerCase()}`;
        const parameterName = 'institution';
        const expectedValue = institutionName.toLowerCase();
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
      });

      test(`returns the correct value for "INSTITUTION=${institutionName}" query parameter`, () => {
          const queryString = `?INSTITUTION=${institutionName}`;
          const parameterName = 'INSTITUTION';
          const expectedValue = institutionName.toLowerCase();
          const returnedValue = getParameterFromQueryString(queryString, parameterName);
          expect(returnedValue).toBe(expectedValue);
      });

      test(`returns the correct value for "INSTITUTION=${institutionName.toLowerCase()}" query parameter`, () => {
        const queryString = `?INSTITUTION=${institutionName.toLowerCase()}`;
        const parameterName = 'INSTITUTION';
        const expectedValue = institutionName.toLowerCase();
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
      });
    }
);

