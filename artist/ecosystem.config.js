module.exports = {
  apps : [{
    name: "artist",
    script: "main.py",
    instances: '1', // Or a number of instances
    interpreter: '/home/artchitector/anaconda3/envs/ldm/bin/python'
  }]
};
