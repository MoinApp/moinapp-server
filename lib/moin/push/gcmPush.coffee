gcm = require 'node-gcm'
uuid = require 'node-uuid'

class GCMPush
  constructor: (gcmApiKey, moinController) ->
    @gcmSender = new gcm.Sender gcmApiKey
    
    moinController.on 'moin', (sender, receipient) =>
      @send sender, receipient
    
  send: (sender, receipient, callback) ->
    if !sender?.getPublicModel?
      return callback new Error 'Must provide database object for "sender".'
    if !receipient?.getPublicModel?
      return callback new Error 'Must provide database object for "receipient".'
    
    receipient.getGcmIDs().complete (err, gcmIDs) =>
      return callback? err if !!err
      return callback? new Error('No device registered for this user.'), true if !gcmIDs || gcmIDs.length == 0
      
      # convert into public model
      gcmIDs = ( gcmID.getPublicModel() for gcmID in gcmIDs )
    
      # now that we have the devices' push tokens
      message = new gcm.Message()
      message.addDataWithObject sender.getPublicModel()
      message.addDataWithKeyValue 'moin-uuid', uuid.v1()
      
      numberOfRetries = 1
      @gcmSender.send message, gcmIDs, numberOfRetries, (err, results) ->
        return callback? err if !!err
        
        if results.success < 1
          text = results.results[0]?.error
          err = new Error "Failed to send GCM push. (\"#{text}\")"
          return callback? err
        
        callback? err, results

module.exports.GCMPush = GCMPush
