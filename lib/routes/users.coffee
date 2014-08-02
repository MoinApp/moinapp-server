db = require '../db/'

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
