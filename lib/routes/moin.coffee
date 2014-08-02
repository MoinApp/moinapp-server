db = require '../db/'

module.exports = (req, res, next) ->
  
  # uid
  to = req.body?.to
  
  if !to
    res.send 400, {
      status: -1,
      error: "Specify id."
    }
  else
    db.User.find({
      where: {
        uid: to
      }
    }).complete (err, user) ->
      throw err if err
      
      if !user
        res.send 400, {
          status: -1,
          error: 'User does not exist.'
        }
      else
        # TODO: MOIN him!
        res.send 200
  
  next()
