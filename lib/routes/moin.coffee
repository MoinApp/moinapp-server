db = require '../db/'
push = require '../push'

_sendPush = (fromUserSessionToken, toUser, callback) ->
  
  db.User.find({
    where: {
      session: fromUserSessionToken
    }
  }).complete (err, fromUser) ->
    callback err
    
    if !fromUser
      return callback new Error 'Sending user not found.'
      
    toUser.getGcmIDs().complete (err, gcmIDs) ->
      callback err
      
      if !gcmIDs
        return callback new Error 'No device registered.'
        
      gcmIDs = ( gcmId.getPublicModel() for gcmId in gcmIDs )
        
      push.sendMessage fromUser, gcmIDs, callback

module.exports = (req, res, next) ->
  
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
        res.send 400, {
          status: -1,
          error: 'User does not exist.'
        }
      else
      
        # TODO: MOIN him!
        res.send 200, {
          status: 0,
          message: "Moin sent."
        }
  
  next()
