import { deepFreeze } from './helpers';

const bannerInstitutionInfo = deepFreeze({
    nyu: {
        logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
        link: 'https://library.nyu.edu',
        imgClass: 'image',
        altLibraryLogoImageText: 'NYU Libraries homepage'
    },
    nyuad: {
        logo: `/images/abudhabi-logo-color.svg`,
        link: 'https://nyuad.nyu.edu/en/library.html',
        imgClass: 'image white-bg',
        altLibraryLogoImageText: 'NYU Abu Dhabi Library homepage'
    },
    nyush: {
        logo: `/images/shanghai-logo-color.svg`,
        link: 'https://shanghai.nyu.edu/academics/library',
        imgClass: 'image white-bg',
        altLibraryLogoImageText: 'NYU Shanghai Library homepage'
    }
});

export { bannerInstitutionInfo };
