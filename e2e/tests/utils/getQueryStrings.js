const fs = require('fs');

const getQueryStrings = () => {
  // Read the JSON file and store the result in a variable
  const TESTCASES = fs.readFileSync(require('path').join(__dirname, '../../../backend/api/testdata/server/test-cases.json'), { encoding: 'utf8'});

  // Parse the JSON data and store the result in a variable
  const jsonData = JSON.parse(TESTCASES);

  const queryStrings = [];

  // Iterate over the JSON object and push the queryString values into the queryStrings array
  jsonData.forEach((item) => {
    queryStrings.push(item.queryString);
  });

  // Return the queryStrings array
  return queryStrings;
};

module.exports = getQueryStrings;