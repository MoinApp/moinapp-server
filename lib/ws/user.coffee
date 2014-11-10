db = require './../db/'

exports.getUser = (username, callback) ->
  if !callback
    # we don't need to perfom a database query if the user wouldn't get the result anyway...
    return
  
  db.User.find({
    where: {
      username: username
    }
  }).complete (err, user) ->
    return callback? err if !!err
    
    if !!user
      user = user.getPublicModel()
      
    callback? err, user
  
exports.findUser = (username, callback) ->
  
  db.User.findAll({
    where: {
      username: {
        like: username + "%"
      }
    }
  }).complete (err, users) ->
    return callback? err if !!err

    publicUsers = []

    users.forEach (user) ->
      publicUsers.push user.getPublicModel()
      
    callback? err, publicUsers
    
exports.getRecents = (user, callback) ->
  
  user.getRecents().complete (err, recents) ->
    return callback? err if !!err
    
    publicRecents = []
    recents.forEach (recent) ->
      publicRecents.push recent.getPublicModel()
      
    callback? err, publicRecents
  