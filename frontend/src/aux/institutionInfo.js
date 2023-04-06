import { deepFreeze } from './helpers';

const bannerInstitutionInfo = deepFreeze({
    nyu: {
        logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
        link: 'https://library.nyu.edu',
        imgClass: 'image',
        altLibraryLogoImageText: 'NYU Libraries homepage.'
    },
    nyuad: {
        logo: 'https://cdn.library.nyu.edu/images/abudhabi_white.svg',
        link: 'https://nyuad.nyu.edu/en/library.html',
        imgClass: 'image',
        altLibraryLogoImageText: 'NYU Abu Dhabi Library homepage.'
    },
    nyush: {
        logo: 'https://cdn.library.nyu.edu/images/shanghai_white.svg',
        link: 'https://shanghai.nyu.edu/academics/library',
        imgClass: 'image',
        altLibraryLogoImageText: 'NYU Shanghai Library homepage.'
    }
});

export { bannerInstitutionInfo };
