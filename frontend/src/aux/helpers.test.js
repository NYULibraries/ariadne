import { getParameterFromQueryString } from './helpers';
import { bannerInstitutionInfo } from './institutionInfo';

const institutionNamesUpperCase = Object.keys(bannerInstitutionInfo).map(institutionName => institutionName.toUpperCase());
describe.each(institutionNamesUpperCase)(
    'Institution name: %s', (institutionNameUpperCase) => {
      test(`returns the correct value for "institution=${institutionNameUpperCase}" query parameter`, () => {
          const queryString = `?institution=${institutionNameUpperCase}`;
          const parameterName = 'institution';
          const expectedValue = institutionNameUpperCase.toLowerCase();
          const returnedValue = getParameterFromQueryString(queryString, parameterName);
          expect(returnedValue).toBe(expectedValue);
      });

      test(`returns the correct value for "institution=${institutionNameUpperCase.toLowerCase()}" query parameter`, () => {
        const queryString = `?institution=${institutionNameUpperCase.toLowerCase()}`;
        const parameterName = 'institution';
        const expectedValue = institutionNameUpperCase.toLowerCase();
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
      });

      test(`returns the correct value for "INSTITUTION=${institutionNameUpperCase}" query parameter`, () => {
          const queryString = `?INSTITUTION=${institutionNameUpperCase}`;
          const parameterName = 'INSTITUTION';
          const expectedValue = institutionNameUpperCase.toLowerCase();
          const returnedValue = getParameterFromQueryString(queryString, parameterName);
          expect(returnedValue).toBe(expectedValue);
      });

      test(`returns the correct value for "INSTITUTION=${institutionNameUpperCase.toLowerCase()}" query parameter`, () => {
        const queryString = `?INSTITUTION=${institutionNameUpperCase.toLowerCase()}`;
        const parameterName = 'INSTITUTION';
        const expectedValue = institutionNameUpperCase.toLowerCase();
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
      });
    }
);

