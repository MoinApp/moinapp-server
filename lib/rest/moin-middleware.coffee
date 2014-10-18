module.exports = (moinController) ->
  # injects the MoinController into the req object
  return (req, res, next) ->
    req.moinController = moinController
    
    next()
