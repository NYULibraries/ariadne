const institutions = {
    nyu: {
        logo: 'https://cdn.library.nyu.edu/images/nyulibraries-logo.svg',
        link: 'http://library.nyu.edu',
        imgClass: 'image'
    },
    nyuad: {
        logo: `${process.env.REACT_APP_PUBLIC_URL}/images/abudhabi-logo-color.svg`,
        link: 'https://nyuad.nyu.edu/en/library.html',
        imgClass: 'image white-bg'
    },
    nyush: {
        logo: `${process.env.REACT_APP_PUBLIC_URL}/images/shanghai-logo-color.svg`,
        link: 'https://shanghai.nyu.edu/academics/library',
        imgClass: 'image white-bg'
    }
};

// eslint-disable-next-line no-console
console.log(institutions.nyuad.logo)
// eslint-disable-next-line no-console
console.log('process.env.PUBLIC_URL:', process.env.REACT_APP_PUBLIC_URL);

export { institutions };

