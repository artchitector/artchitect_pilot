module.exports = {
  apps : [{
    name: "artist",
    script: "main.py",
    instances: '1', // Or a number of instances
    interpreter: '/home/artchitector/anaconda3/envs/artchitect/bin/python',
    env: {
      INVOKEAI_ROOT: "/home/artchitector/invoke-ai/invokeai_v2.2.5/invokeai"
    },
  }]
};
