restify = require 'restify'
db = require '../../db/'

exports.GETuser = (req, res, next) ->
  
  username = req.params?.username
  
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

exports.POSTaddGcm = (req, res, next) ->
  
  gcmId = req.body?.gcmId
  
  db.gcmID.find({
    where: {
      uid: gcmId
    }
  }).complete (err, id) ->
    return next err if !!err
    if !!id
      return next new restify.InvalidArgumentError 'GCM ID is already added.'
    
    req.user.addGcmID(gcmId).complete (err) ->
      return next err if !!err
      
      next()
