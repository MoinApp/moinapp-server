db = require './../db/'

exports.getUser = (username, callback) ->
  db.User.find({
    where: {
    }
  }).complete (err, user) ->
    return callback? err if !!err
    
    if !!user
      user = user.getPublicModel()
      
    callback? err, user
  