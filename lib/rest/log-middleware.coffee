module.exports = (req, res, next) ->
  
  ip = req.connection.remoteAddress
  method = req.route.method
  path = req.path()
  text = ip + ": " + method + " " + path
  
  console.log text
  
  next()