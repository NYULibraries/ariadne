const getLinks = (jsonData) => {
  const result = jsonData.ctx_obj[0].ctx_obj_targets[0].target;
  return result;
};

const getCoverageStatement = (link) => {
  return link.coverage[0].coverage_text[0].threshold_text[0].coverage_statement.join('. ');
};

export { getLinks, getCoverageStatement };
