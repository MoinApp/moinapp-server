restify = require 'restify'
db = require './../../db/'
{ SessionHandler } = require './../../auth/sessionHandler'

exports.checkAuthentication = (req, res, next) ->

  user = req.user
  if !!user
    next()
  else
    sessionToken = req.headers?.authorization

    if !sessionToken
      next new restify.NotAuthorizedError 'Requires session token.'
    else
      SessionHandler.getInstance().getUserForSessionToken sessionToken, (err, user) ->
        return next err if !!err

        req.user = user

        next()

exports.POSTsignin = (req, res, next) ->

  username = req.body?.username.trim()
  password = req.body?.password.trim()
  app = req.headers?["user-agent"]

  if !username || !password || !app
    return next new restify.InvalidArgumentError 'Please provide all login information.'

  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err if !!err

    if !user
      next new restify.InvalidCredentialsError 'Username not found.'
    else

      if !user.isValidPassword password
        return next new restify.InvalidCredentialsError 'Invalid password.'

      db.Session.createNew(app).complete (err, session) ->
        return next err if !!err

        user.addSession session

        res.send 200, session.getPublicModel()

        next()

exports.POSTsignup = (req, res, next) ->

  username = req.body?.username.trim()
  password = req.body?.password.trim()
  app = req.headers?["user-agent"]
  # optional, but clients should send this
  email = req.body?.email.trim()

  if !username || !password || !app
    return next new restify.InvalidArgumentError 'Please provide username and password at least.'

  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err if !!err

    if !!user
      next new restify.InvalidArgumentError 'Username is already taken.'
    else
      db.User.createUserAndEncryptPassword({
        username: username,
        password: password,
        email: email
      }).complete (err, user) ->
        return next err if !!err

        db.Session.createNew(app).complete (err, session) ->
          return next err if !!err

          user.addSession session

          res.send 200, session.getPublicModel()
          next()
