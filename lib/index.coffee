cluster = require 'cluster'

if cluster.isMaster
  # first start
  console.log "Moin"
  
  os = require 'os'
  numCPUs = os.cpus().length
  console.log "Cluster:", "Could spawn #{numCPUs} forks."
  
  numberOfForks = 0
  
  startFork = ->
    cluster.fork()
  
  cluster.on 'listening', (worker) ->
    numberOfForks++
    console.log "Currently running #{numberOfForks} fork(s)."
    
  cluster.on 'exit', (worker, code, signal) ->
    numberOfForks--
    
    console.log "Worker #{worker} died with signal #{signal}."
    console.log "Currently running #{numberOfForks} fork(s)."
    
    if signal == null
      console.log "Expecting syntax error. Stopping."
    else
      startFork()
      console.log "Reforked."
  
  # fork 1 for now
  startFork()
  # TODO: if we want to use clusters, better only use db.coffee only once
  # toughts: maybe just instantiate once in the master and communicate with it from the forks via messages to get db objects?
  #for i in [1...numCPUs] by 1
  #  startFork()

else
  # x'th start
  
  require('./fork')()
