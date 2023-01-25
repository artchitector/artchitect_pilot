module.exports = {
  apps : [{
    name: "origin",
    script: "./bin/origin",
    instances: '1',
    args: "-l 0.0.0.0:8081 -s 640x480"
  }]
};
