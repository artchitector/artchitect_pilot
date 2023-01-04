module.exports = {
  apps: [
    {
      name: 'Artchitect Origin',
      exec_mode: 'cluster',
      instances: '1', // Or a number of instances
      script: './main.py',
      args: 'start'
    }
  ]
}
