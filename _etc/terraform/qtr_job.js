const { execSync } = require("child_process");

exports.handler = async (event) => {
  console.log(`run qtr-job: event=${JSON.stringify(event)}`);

  // NOTE: for debug (make a function fail)
  if (event["_fail"]) {
    throw "an error occurred";
  }

  // do something
  const zen =
    "Zen: " + execSync("curl -sf https://api.github.com/zen").toString();
  console.log(`********** ${zen} **********`);
};
