module.exports = {
  apps : [{
    name: "origin",
    script: "main.py",
    instances: '1', // Or a number of instances
    interpreter: '/home/artchitector/anaconda3/envs/artchitect/bin/python'
  }]
};
