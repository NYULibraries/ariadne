import {
  NYUAD,
  NYUSH,
  DEFAULT_LOGO,
  DEFAULT_LINK,
  DEFAULT_IMG_CLASS,
  NYUAD_LOGO,
  NYUAD_LINK,
  NYUAD_IMG_CLASS,
  NYUSH_LOGO,
  NYUSH_LINK,
  NYUSH_IMG_CLASS,
} from './institutionConstants';

//  Helper functions for useApi hook
const getLinks = (jsonData) => {
  return jsonData.ctx_obj[0].ctx_obj_targets[0].target;
};

// Helper function for List.js component
const getCoverageStatement = (link) => {
  return link.coverage[0].coverage_text[0].threshold_text[0].coverage_statement?.join('. ');
};

// Helper functions for Banner.js component
const getParameterFromQueryString = (queryString, parameterName) => {
  const urlParams = new URLSearchParams(queryString);
  let parameter = urlParams.get(parameterName);
  if (urlParams.has(`umlaut.${parameterName}`)) {
    parameter = urlParams.get(`umlaut.${parameterName}`);
    urlParams.delete(`umlaut.${parameterName}`);
    urlParams.set(parameterName, parameter);
  }
  return parameter;
};

const getQueryParameter = (parameterName) => {
  const queryString = window.location.search;
  return getParameterFromQueryString(queryString, parameterName);
};

// Helper function for Banner.js component

const getInstitution = (institution) => {
  let logo = DEFAULT_LOGO;
  let link = DEFAULT_LINK;
  let imgClass = DEFAULT_IMG_CLASS;

  if (institution === NYUAD) {
    logo = NYUAD_LOGO;
    link = NYUAD_LINK;
    imgClass = NYUAD_IMG_CLASS;
  } else if (institution === NYUSH) {
    logo = NYUSH_LOGO;
    link = NYUSH_LINK;
    imgClass = NYUSH_IMG_CLASS;
  }

  return { logo, link, imgClass };
};

export { getLinks, getCoverageStatement, getQueryParameter, getInstitution, getParameterFromQueryString };
