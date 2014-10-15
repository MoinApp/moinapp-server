db = require '../../db/'
{ MoinController } = require '../../moin/moinController'

exports.POSTmoin = (req, res, next) ->
  
  sender = req.user?.username
  receipient = req.body?.username
  
  moin = new MoinController
  moin.setUsersFromNames sender, receipient
  moin.sendMoin (err) ->
    next err
