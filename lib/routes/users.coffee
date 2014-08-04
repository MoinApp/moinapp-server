restify = require 'restify'
db = require '../db/'
session = require './session'
crypt = require '../db/crypt'
util = require 'util'

class UsernameTakenError extends restify.RestError
  constructor: (@message) ->
    restify.RestError.call this, {
      restCode: 'UsernameTakenError',
      statusCode: 409,
      message: message,
      constructorOpt: UsernameTakenError
    }
    @name = 'UsernameTakenError'
  

exports.getUser = (req, res, next) ->
  
  username = req.params?.username
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err if !!err
    
    if !user
      res.send 404, {}
    else
      res.send user.getPublicModel()
      
    next()

exports.newUser = (req, res, next) ->
  
  username = req.body?.username
  password = req.body?.password
  email = req.body?.email
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err if !!err
    
    if !!user
      next new UsernameTakenError 'Username is already taken.'
    else
      
      if username.length < 3
        return next new restify.InvalidArgumentError 'Username is too short.'
      if password.length < 5
        return next new restify.InvalidArgumentError 'Password is too short.'
      
      db.User.createUser({
        username: username,
        password: crypt.getSHA256(password),
        emailHash: crypt.getMD5(email)
      }).complete (err, user) ->
        return next err if !!err
        
        session.createSession user.username, (err, sessionToken) ->
          return next err if !!err
          
          res.send {
            code: "Success",
            session: sessionToken
          }
          
          next()
          
exports.signIn = (req, res, next) ->
  
  username = req.body?.username
  password = req.body?.password
  
  db.User.find({
    where: {
      username: username,
      password: crypt.getSHA256 password
    }
  }).complete (err, user) ->
    return next err if !!err
    
    if !user
      res.send 400, {
        code: "CredentialsWrong",
        message: 'Username or Password is wrong.'
      }
      next()
    else
      session.getOrCreateSession username, (err, sessionToken) ->
        return next err if !!err
        
        res.send 200, {
          code: "Success",
          session: sessionToken
        }
        
        next()
        
exports.addGCMId = (req, res, next) ->
  
  sessionToken = req.query?.session
  gcmIdString = req.body?.gcm_id
  
  db.User.find({
    where: {
      session: sessionToken
    }
  }).complete (err, user) ->
    return next err if !!err
    
    if !!user
      db.gcmID.findOrCreate({
        uid: gcmIdString
      }).complete (err, gcmId, created) ->
        return next err if !!err
        
        if created
          return next new restify.InvalidArgumentError 'gcmId already exists.'
        
        user.addGcmID(gcmId).complete (err) ->
          return next err if !!err
          
          res.send 200, {
            code: "GCMIDAdded",
            message: "GcmID added."
          }
          
          next()
    else
      next new restify.InvalidCredentialsError 'User could not be found for session.'
        
exports.validateSession = (req, res, next) ->
  
  sessionToken = req.body?.session
  
  session.validateSession sessionToken, (err, ok) ->
    throw err if !!err
    
    if ok
      res.send 200, {
        status: 0,
        message: "Session valid."
      }
    else
      res.send 400, {
        status: -1,
        message: "Invalid session."
      }
      
    next()
