uuid = require 'node-uuid'
db = require '../db/'

exports.getSession = (username, callback) ->
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    throw err if err
    
    if !user
      callback new Error('User not found.')
    else
      callback null, user.session

exports.createSession = (username, callback) ->
    
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    throw err if err
    
    if !user
      callback new Error('User not found.')
    else
      
      uid = uuid.v1() # time-based
      user.session = uid
      
      user.save().complete (err) ->
        throw err if err
        
        callback null, uid
  
exports.getOrCreateSession = (username, callback) ->
  
  exports.getSession username, (err, sessionToken) ->
    throw err if err
    
    if sessionToken
      callback null, sessionToken
    else
      
      exports.createSession username, (err, sessionToken) ->
        callback err, sessionToken
