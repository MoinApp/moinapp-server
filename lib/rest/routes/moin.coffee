db = require '../../db/'

exports.POSTmoin = (req, res, next) ->
  
  sender = req.user?.username
  receipient = req.body?.username
  
  req.moinController.sendMoin sender, receipient, (err, warnings) ->
    return next err if !!err
    
    res.send 200, {
      code: "Success",
      message: warnings
    }
    next()
