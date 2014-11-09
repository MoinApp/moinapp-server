pkg = require '../../../package'

exports.GETindex = (req, res, next) ->
  # redirect to the github homepage
  res.setHeader 'Location', "http://i.imgur.com/E2T98iu.jpg" #pkg.homepage
  res.send 302
  next()
