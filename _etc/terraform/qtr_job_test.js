exports.handler = async (event) => {
  if (event["_fail"]) {
    throw "error";
  }
};
