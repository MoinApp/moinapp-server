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
  if !sendingUser?.getPublicModel?
    return callback new Error 'Must provide database object for sending user.'
  if toGCMIDs.length == 0
    return callback new Error 'No receiver GCM IDs given.'
  
  exports._checkInit()
  
  message = new gcm.Message()
  #message.addDataWithKeyValue 'from', sendingUser.getPublicModel()
  message.addDataWithObject sendingUser.getPublicModel()
  
  numberOfRetries = 1
  
  #console.log "Sending message from #{sendingUser.username} to" , toGCMIDs, "..."
  gcmSender?.send message, toGCMIDs, numberOfRetries, callback
