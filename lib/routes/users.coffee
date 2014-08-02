db = require '../db/'
session = require './session'

exports.getUser = (req, res, next) ->
  
  username = req.params?.username
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    throw err if err
    
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
    throw err if err
    
    if !!user
      res.send 400, { status: -1, error: 'Username is already taken.' }
    else
      
      db.User.createUser({
        username: username,
        password: password
      }).complete (err, user) ->
        throw err if err
        
        session.createSession user.username, (err, sessionToken) ->
          throw err if err
          
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
    throw err if err
    
    if !user
      res.send 400, {
        status: -1,
        error: 'Username or Password is wrong.'
      }
    else
      session.getOrCreateSession username, (err, sessionToken) ->
        throw err if err
        
        res.send 200, {
          status: 0,
          session: sessionToken
        }
