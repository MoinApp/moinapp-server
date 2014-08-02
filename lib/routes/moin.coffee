restify = require 'restify'
db = require '../db/'
push = require '../push'

exports._sendPush = (fromUserSessionToken, toUser, callback) ->
  
  db.User.find({
    where: {
      session: fromUserSessionToken
    }
  }).complete (err, fromUser) ->
    callback err
    
    if !fromUser
      return callback new Error 'Sending user not found.'
      
    toUser.getGcmIDs().complete (err, gcmIDs) ->
      return callback err if !!err
      
      if !gcmIDs
        return callback new Error 'No device registered.'
        
      gcmIDs = ( gcmId.getPublicModel() for gcmId in gcmIDs )
      
      console.log "Sending moin from #{fromUser.username} to #{toUser.username}..."
      push.sendMessage fromUser, gcmIDs, callback

exports.moin = (req, res, next) ->
  
  # uid
  to = req.body?.to
  
  if !to
    res.send 400, {
      status: -1,
      error: "Specify id."
    }
  else
    db.User.find({
      where: {
        uid: to
      }
    }).complete (err, user) ->
      return next err if !!err
      
      if !user
        next new restify.ResourceNotFoundError 'User does not exist.'
      else
        # Send Moin!
        from_sessionToken = req.query?.session
        exports._sendPush from_sessionToken, user, (err, results) ->
          return next err if !!err
          
          console.log "Result from gcm send:", results
          
          res.send 200, {
            status: 0,
            message: "Moin sent."
          }
