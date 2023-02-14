function updateGoldenFiles() {
  return process.env.UPDATE_GOLDEN_FILES &&
         process.env.UPDATE_GOLDEN_FILES.toLowerCase() !== 'false';
}

export {
  updateGoldenFiles,
};
