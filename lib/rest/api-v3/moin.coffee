restify = require 'restify'
db = require '../../db/'

exports.POSTmoin = (req, res, next) ->

  sender = req.user?.username
  recipient = req.body?.username.trim()

  if !sender
    return next new restify.InternalError 'You are not authorized.'
  if !sender || !recipient
    return next new restify.InvalidArgumentError 'Specify a recipient.'

  req.moinController.sendMoin sender, recipient, (err) ->
    return next err if !!err

    res.send 200
    next()
