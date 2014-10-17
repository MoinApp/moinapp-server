db = require '../../db/'
{ MoinController } = require '../../moin/moinController'

exports.POSTmoin = (req, res, next) ->
  
  sender = req.user?.username
  receipient = req.body?.username
  
  moin = new MoinController
  moin.setUsersFromNames sender, receipient, (err) ->
    return next err if !!err
    
    # done setting the user objects
    moin.sendMoin (err) ->
      next err
