import { deepFreeze } from './helpers';

const bannerInstitutionInfo = deepFreeze({
    nyu: {
        logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
        link: 'http://library.nyu.edu',
        imgClass: 'image'
    },
    nyuad: {
        logo: `/images/abudhabi-logo-color.svg`,
        link: 'https://nyuad.nyu.edu/en/library.html',
        imgClass: 'image white-bg'
    },
    nyush: {
        logo: `/images/shanghai-logo-color.svg`,
        link: 'https://shanghai.nyu.edu/academics/library',
        imgClass: 'image white-bg'
    }
});

export { bannerInstitutionInfo };

