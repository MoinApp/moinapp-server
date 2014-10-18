db = require './../db/'

exports.getUser = (username, callback) ->
  if !callback
    # we don't need to perfom a database query if the user wouldn't get the result anyway...
    return
  
  db.User.find({
    where: {
    }
  }).complete (err, user) ->
    return callback? err if !!err
    
    if !!user
      user = user.getPublicModel()
      
    callback? err, user
  