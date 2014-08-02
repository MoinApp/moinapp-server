cluster = require 'cluster'

if cluster.isMaster
  # first start
  console.log "Moin"
  
  os = require 'os'
  numCPUs = os.cpus().length
  console.log "Cluster:", "Could spawn #{numCPUs} forks."
  
  startFork = ->
    cluster.fork()
  
  cluster.on 'exit', (worker, code, signal) ->
    console.log "Worker #{worker} died with signal #{signal}."
    
    if signal == null
      console.log "Expecting syntax error. Stopping."
    else
      startFork()
      console.log "Reforked."
  # fork 1 for now
  startFork()

else
  # x'th start
  
  require('./fork')()
