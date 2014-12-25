fs = require 'fs'
{ EventEmitter } = require 'events'
db = require '../db/'
{ GCMPush } = require './push/gcmPush'
{ APNPush } = require './push/apnPush'

class MoinController extends EventEmitter
  constructor: ->
    @androidPush = new GCMPush process.env.GCM_API_KEY, this
    @iOSPush = new APNPush @getAPNCertificate(), this

  getAPNCertificate: ->
    certString = process.env.APN_CERT
    certFilename = process.env.APN_CERT_FILE
    if certString
      new Buffer(certString, 'binary')
    else if certFilename
      fs.readFileSync certFilename
    else
      new Buffer(0)
    
  _getUsersFromNames: (senderName, receipientName, callback) ->
    @_resolveUser senderName, (err, sender) =>
      return callback err if !!err
      return callback new Error 'User "' + senderName + '" not found.' if !sender

      @_resolveUser receipientName, (err, receipient) =>
        return callback err if !!err
        return callback new Error 'User "' + receipientName + '" not found.' if !receipient

        callback null, sender, receipient

  _resolveUser: (username, callback) ->
    db.User.find({
      where: {
        username
      }
    }).complete callback

  _addReceivingUserToRecents: (sender, recipient, callback) ->
    sender.removeRecent(recipient).complete (err) ->
      # adds the receiving user to the "recents" list of the sender
      sender.addRecent(recipient).complete callback

  _sendMoinEvent: (sender, receipient) ->
    @emit 'moin', sender, receipient
  sendMoin: (senderName, receipientName, callback) ->
    @_getUsersFromNames senderName, receipientName, (err, sender, receipient) =>
      return callback? err if !!err

      @_addReceivingUserToRecents sender, receipient, (err) =>
        return callback? err if !!err

        # users are validated at this point
        # send event
        @_sendMoinEvent sender, receipient

        callback? null, null

module.exports.MoinController = MoinController
