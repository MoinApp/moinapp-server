module.exports = (moinController) ->
  return (req, res, next) ->
    req.moinController = moinController
    
    next()
