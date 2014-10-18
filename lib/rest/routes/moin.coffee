db = require '../../db/'

exports.POSTmoin = (req, res, next) ->
  
  sender = req.user?.username
  receipient = req.body?.username
  
  req.moinController.sendMoin sender, receipient, (err) ->
    return next err if !!err
    
    # done setting the user objects
    moin.sendMoin (err, warnings) ->
      return next err if !!err
      
      res.send 200, {
        code: "Success",
        message: warnings
      }
      next()
