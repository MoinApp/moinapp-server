db = require '../db/'
session = require './session'

exports.getUser = (req, res, next) ->
  
  username = req.params?.username
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err
    
    if !user
      res.send 404, {}
    else
      res.send user.getPublicModel()
      
    next()

exports.newUser = (req, res, next) ->
  
  username = req.body?.username
  password = req.body?.password
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return next err
    
    if !!user
      res.send 400, { status: -1, error: 'Username is already taken.' }
    else
      
      db.User.createUser({
        username: username,
        password: password
      }).complete (err, user) ->
        return next err
        
        session.createSession user.username, (err, sessionToken) ->
          return next err
          
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
    return next err
    
    if !user
      res.send 400, {
        status: -1,
        error: 'Username or Password is wrong.'
      }
      next()
    else
      session.getOrCreateSession username, (err, sessionToken) ->
        return next err
        
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
    return next err
    
    if !!user
      db.gcmID.create({
        uid: gcmIdString
      }).complete (err, gcmId) ->
        return next err
        
        user.addGcmID(gcmId).complete (err) ->
          return next err
          
          res.send 200, {
            status: 0,
            message: "GcmID added."
          }
          
          next()
        
