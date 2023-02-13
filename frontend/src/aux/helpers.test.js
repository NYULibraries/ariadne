import { getParameterFromQueryString } from './helpers';

describe('getParameterFromQueryString', () => {
    test('returns the correct value for "institution" query parameter', () => {
        const queryString = '?institution=NYU';
        const parameterName = 'institution';
        const expectedValue = 'NYU';
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
    });

    test('returns the correct value for "umlaut.institution" query parameter', () => {
        const queryString = '?umlaut.institution=NYUAD';
        const parameterName = 'institution';
        const expectedValue = 'NYUAD';
        const returnedValue = getParameterFromQueryString(queryString, parameterName);
        expect(returnedValue).toBe(expectedValue);
    });
});
