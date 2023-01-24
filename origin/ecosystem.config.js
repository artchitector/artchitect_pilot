module.exports = {
  apps : [{
    name: "origin",
    script: "./bin/origin",
    instances: '1',
    args: "-l localhost:8081 -s 640x480"
  }]
};
