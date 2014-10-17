restify = require 'restify'
db = require './../../db/'

class SessionHandler
  @instance = null
  @getInstance = ->
    if !SessionHandler.instance
      SessionHandler.instance = new SessionHandler
    SessionHandler.instance
  
  constructor: ->
    
  getUserForSessionToken: (token, callback) ->
    db.Session.find({
      where: {
        uid: token
      }
    }).complete (err, session) ->
      return callback err if !!err
      
      console.log "session:", session
      if !session
        callback new restify.NotAuthorizedError 'Invalid session token.'
      else
        session.getUser().complete callback
    
exports.checkAuthentication = (req, res, next) ->
  
  user = req.user
  if !!user
    next()
  else
    sessionToken = req.query?.session
  
    if !sessionToken
      next new restify.NotAuthorizedError 'Requires session token.'
    else
      SessionHandler.getInstance().getUserForSessionToken sessionToken, (err, user) ->
        return next err if !!err
        
        req.user = user
        
        next()
  
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
        return next new restify.InvalidCredentialsError 'Invalid password.'
        
      db.Session.createNew(app).complete (err, session) ->
        return next err if !!err
        
        user.addSession session
        
        req.user = user
        
        res.send 200, {
          code: "Success",
          message: session.getPublicModel()
        }
        
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
          req.user = user
        
          next()
