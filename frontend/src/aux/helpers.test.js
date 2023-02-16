import { getParameterFromQueryString } from './helpers';

describe('getParameterFromQueryString', () => {
    test('returns the correct value for "institution=NYU" query parameter', () => {
        const queryString = '?institution=NYU';
        const parameterName = 'institution';
        const expectedValue = 'NYU';
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
    });

    test('returns the correct value for "institution=NYUAD" query parameter', () => {
        const queryString = '?institution=NYUAD';
        const parameterName = 'institution';
        const expectedValue = 'NYUAD';
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
    });

    test('returns the correct value for "institution=NYUSH" query parameter', () => {
        const queryString = '?institution=NYUSH';
        const parameterName = 'institution';
        const expectedValue = 'NYUSH';
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
    });

});
