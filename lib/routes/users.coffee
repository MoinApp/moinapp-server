restify = require 'restify'
db = require '../db/'
session = require './session'
crypt = require '../db/crypt'

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
      next new restify.RestError 'Username is already taken.'
    else
      
      if username.length < 3
        return next new restify.InvalidArgumentError 'Username is too short.'
      if password.length < 5
        return next new restify.InvalidArgumentError 'Password is too short.'
      
      db.User.createUser({
        username: username,
        password: password,
        emailHash: crypt.getMD5 email
      }).complete (err, user) ->
        return next err if !!err
        
        session.createSession user.username, (err, sessionToken) ->
          return next err if !!err
          
          res.send {
            status: 0,
            session: sessionToken
          }
          
          next()
          
exports.signIn = (req, res, next) ->
  
  username = req.body?.username
  password = req.body?.password
  
  db.User.find({
    where: {
      username: username,
      password: password
    }
  }).complete (err, user) ->
    return next err if !!err
    
    if !user
      res.send 400, {
        status: -1,
        error: 'Username or Password is wrong.'
      }
      next()
    else
      session.getOrCreateSession username, (err, sessionToken) ->
        return next err if !!err
        
        res.send 200, {
          status: 0,
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
      db.gcmID.find({
        where: {
          uid: gcmIdString
        }
      }).complete (err, existingGcmID) ->
        return next err if !!err
        
        if !!existingGcmID
          return next new restify.InvalidArgument 'GcmID already existing.'
          
        db.gcmID.create({
          uid: gcmIdString
        }).complete (err, gcmId) ->
          return next err if !!err
          
          user.addGcmID(gcmId).complete (err) ->
            return next err if !!err
            
            res.send 200, {
              status: 0,
              message: "GcmID added."
            }
            
            next()
    else
      next new restify.InvalidCredentialsError 'User could not be found for session.'
        
