restify = require 'restify'
db = require './../../db/'

class SessionHandler
  @instance = null
  @getInstance = ->
    if !@instance
      @instance = new SessionHandler
    @instance
  
  constructor: ->
    
  checkSessionToken: (token) ->
    !!token
    
exports.checkAuthentication = (req, res, next) ->
  
  sessionToken = req.query?.session
  
  if SessionHandler.getInstance().checkSessionToken sessionToken
    next()
  else
    next new restify.NotAuthorizedError 'Requires session token.'
  
exports.POSTsignin = (req, res, next) ->
  
  username = req.body?.username
  password = req.body?.password
  app = req.body?.application
  
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
        return new restify.InvalidCredentialsError 'Invalid password.'
        
      db.Session.createNew(app).complete (err, session) ->
        return next err if !!err
        
        user.addSession session
        
        req.session = session
        req.user = user
        
        next()

exports.POSTsignup = (req, res, next) ->
  
  username = req.body?.username
  password = req.body?.username
  app = req.body?.application
  # optional
  email = req.body?.email
  
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
          req.session = session
          req.user = user
        
          next()
