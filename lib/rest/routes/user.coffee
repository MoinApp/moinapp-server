restify = require 'restify'
db = require '../../db/'

exports.GETuser = (req, res, next) ->
  
  username = req.param?.username
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err if !!err
    
    if !user
      next new restify.ResourceNotFoundError 'User does not exist.'
    else
      res.send 200, user.getPublicModel()
      
      next()
