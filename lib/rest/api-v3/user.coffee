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
      res.send 200, {
        code: "Success",
        message: user.getPublicModel()
      }

      next()

exports.GETuserSearch = (req, res, next) ->

  username = req.params?.username

  db.User.findAll({
    where: {
      username: {
        like: "%" + username + "%"
      }
    }
  }).complete (err, users) ->
    return next err if !!err

    publicUsers = []

    users.forEach (user) ->
      publicUsers.push user.getPublicModel()

    res.send 200, {
      code: "Success",
      message: publicUsers
    }
    
    next()
    
exports.GETusersRecents = (req, res, next) ->
  
  req.user.getRecents().complete (err, recents) ->
    return next err if !!err
    
    publicRecents = []
    recents.forEach (recent) ->
      publicRecents.push recent.getPublicModel()
      
    res.send 200, {
      code: "Success",
      message: publicRecents
    }

exports.POSTaddGcm = (req, res, next) ->

  gcmId = req.body?.gcmId

  if !gcmId
    return next new restify.InvalidArgumentError 'Specify a GCM ID.'

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

      res.send 200, {
        code: "Success",
        message: "GCM ID added."
      }

      next()
