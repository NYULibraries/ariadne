import { deepFreeze } from './helpers';

const bannerInstitutionInfo = deepFreeze({
    nyu: {
        logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
        link: 'http://library.nyu.edu',
        imgClass: 'image',
        altLibraryLogoImageText: 'NYU Libraries logo - click to go to the NYU Libraries main page'
    },
    nyuad: {
        logo: `/images/abudhabi-logo-color.svg`,
        link: 'https://nyuad.nyu.edu/en/library.html',
        imgClass: 'image white-bg',
        altLibraryLogoImageText: 'NYU Abu Dhabi Library logo - click to go to the NYU Abu Dhabi Library main page'
    },
    nyush: {
        logo: `/images/shanghai-logo-color.svg`,
        link: 'https://shanghai.nyu.edu/academics/library',
        imgClass: 'image white-bg',
        altLibraryLogoImageText: 'NYU Shanghai Library logo - click to go to the NYU Shanghai Library main page'
    }
});

export { bannerInstitutionInfo };
