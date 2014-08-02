db = require '../db/'

exports.getUser = (req, res, next) ->
  
  name = req.params?.name
  
  db.User.find({
    where: {
      username: name
    }
  }).complete (err, user) ->
    throw err if err
    
    if !user
      res.send 404, {}
    else
      res.send user.getPublicModel()
      
    next()
