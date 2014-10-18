db = require '../../db/'

exports.POSTmoin = (req, res, next) ->
  
  sender = req.user?.username
  receipient = req.body?.username
  
  if !sender
    return next new restify.InternalError 'You are not authorized.'
  if !sender || !receipient
    return next new restify.InvalidArgumentError 'Specify a receipient.'
  
  req.moinController.sendMoin sender, receipient, (err, warnings) ->
    return next err if !!err
    
    res.send 200, {
      code: "Success",
      message: warnings
    }
    next()
