module.exports = (req, res, next) ->
  
  method = req.route.method
  path = req.route.path
  text = method + " " + path
  
  console.log text
  
  next()