gcm = require 'node-gcm'

gcmSender = null

exports._checkInit = ->
  if !gcmSender
    apiKey = process.env.GCM_API_KEY
    
    if !apiKey
      throw new Error 'GCM_API_KEY not found.'
    exports._init apiKey
    
exports._init = (apiKey) ->
  gcmSender = new gcm.Sender apiKey

exports.sendMessage = (sendingUser, toGCMIDs, callback) ->
  exports._checkInit()
  
  message = new gcm.Message
  #message.addDataWithKeyValue 'from', sendingUser.getPublicModel()
  message.addDataWithObject sendingUser.getPublicModel()
  
  numberOfRetries = 1
  
  gcmSender?.send message, toGCMIDs, numberOfRetries, (err, result) ->
    console.log "Result from gcm send:", result
    callback err, result
