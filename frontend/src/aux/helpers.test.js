import { getInstitution, getParameterFromQueryString } from './helpers';

describe('getInstitution', () => {
    test('returns the correct object for institution "NYUAD"', () => {
        const institution = 'NYUAD';
        const result = getInstitution(institution);
        expect(result).toEqual({
            logo: `${process.env.PUBLIC_URL}/images/abudhabi-logo-color.svg`,
            link: 'https://nyuad.nyu.edu/en/library.html',
            imgClass: 'image white-bg',
        });
    });
    test('returns the correct object for institution "NYUSH"', () => {
        const institution = 'NYUSH';
        const result = getInstitution(institution);
        expect(result).toEqual({
            logo: `${process.env.PUBLIC_URL}/images/shanghai-logo-color.svg`,
            link: 'https://shanghai.nyu.edu/academics/library',
            imgClass: 'image white-bg',
        });
    });
    test('returns the default object for NYU institution', () => {
        const institution = 'NYU';
        const result = getInstitution(institution);
        expect(result).toEqual({
            logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
            link: 'http://library.nyu.edu',
            imgClass: 'image',
        });
    });

    test('returns the default object for unknown institution', () => {
        const institution = 'unknown';
        const result = getInstitution(institution);
        expect(result).toEqual({
            logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
            link: 'http://library.nyu.edu',
            imgClass: 'image',
        });
    });
});

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
