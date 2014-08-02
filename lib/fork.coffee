server = require './server'

module.exports = () ->
  console.log "Fork started."
  server.init()
  server.start()
