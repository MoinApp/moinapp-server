gcm = require 'node-gcm'
uuid = require 'node-uuid'

gcmSender = null

exports._checkInit = (callback) ->
  if !gcmSender
    apiKey = process.env.GCM_API_KEY
    
    if !apiKey
      err = new Error 'GCM_API_KEY not found.'
      if process.env.NODE_ENV == 'production'
        console.error err.message
      else
        console.warn err.message
        
      callback err
    else
      exports._init apiKey
      callback null
    
exports._init = (apiKey) ->
  gcmSender = new gcm.Sender apiKey
  return true

exports.sendMessage = (sendingUser, toGCMIDs, callback) ->
  if !sendingUser?.getPublicModel?
    return callback new Error 'Must provide database object for sending user.'
  if toGCMIDs.length == 0
    return callback new Error 'No receiver GCM IDs given.'
  
  exports._checkInit (err) ->
    if !!err
      return callback err
    
    message = new gcm.Message()
    #message.addDataWithKeyValue 'from', sendingUser.getPublicModel()
    message.addDataWithObject sendingUser.getPublicModel()
    message.addDataWithKeyValue 'moin-uuid', uuid.v1()
    
    numberOfRetries = 1
    
    #console.log "Sending message from #{sendingUser.username} to" , toGCMIDs, "..."
    gcmSender?.send message, toGCMIDs, numberOfRetries, callback
