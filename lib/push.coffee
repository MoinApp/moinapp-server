# TODO use GCM to push moins
gcm = require 'node-gcm'

gcmSender = null

exports._checkInit = ->
  if !gcmSender
    exports._init process.env.GCM_API_KEY
    
exports._init = (apiKey) ->
  gcmSender = new gcm.Sender apiKey

exports.sendMessage = (sendingUser, toGCMIDs) ->
  exports._checkInit()
  
  message = new gcm.Message
  #message.addDataWithKeyValue 'from', sendingUser.getPublicModel()
  message.addDataWithObject sendingUser.getPublicModel()
  
  numberOfRetries = 1
  
  gcmSender?.send message, toGCMIDs, numberOfRetries, (err, result) ->
    console.log "Result from gcm send:", result
