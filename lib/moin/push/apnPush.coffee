apn = require 'apn'
db = require '../../db/'

class APNPush
  constructor: (pfxBuffer, moinController) ->
    # do some init here
    if !pfxBuffer || pfxBuffer.length == 0
      console.log "Warning: No Apple Push Notification certificate provided."
      return

    # see https://github.com/argon/node-apn/blob/master/doc/connection.markdown
    @apnConnection = new apn.Connection {
      pfx: pfxBuffer,
      connectionTimeout: 5000
    }
    @apnConnection.on 'error', (err) =>
      console.log "APN connection error:", err
      @isRunning = false
    @apnConnection.on 'transmissionError', @handleError

    @feedback = new apn.Feedback {
      "interval": ( 6 * 60 * 60 ) # check every 6 hrs
    }
    @feedback.on 'feedback', (time, deviceTokenBuffer) =>
      deviceToken = deviceTokenBuffer.toString()
      @deleteDeviceToken deviceToken
    console.log "APN feedback service started."

    moinController.on 'moin', (sender, receipient) =>
      @send sender, receipient
    console.log "APN Push running."
    @isRunning = true

  handleError: (errorCode, notification, device) ->
    errorMessage = errorCode
    for key of apn.Errors
      if apn.Errors[key] == errorCode
        errorMessage = key
        break

    console.log "APN transmission error:", errorMessage + " (Code: " + errorCode + ")", "for notification:", notification, "to device:", device

    if errorCode == 8 # invalid token
      token = device.token.toString 'hex'
      console.log "Searching for '" + token + "'..."
      db.APNDeviceToken.find({
        where: {
          uid: token
        }
      }).complete (err, deviceToken) ->
        if !!err
          console.log err

        deviceToken.destroy().complete (err) ->
          if !!err
            console.log err
          else
            console.log "Device token removed."

  deleteDeviceToken: (deviceToken) ->
    # TODO: implement
    db.APNDeviceToken.find({
      where: {
        uid: deviceToken
      }
    }).complete (err, db_deviceToken) ->
      if ( !!err )
        console.log "APN feedback: Deleting '" + deviceToken + "' errored: " + err
        return

      db_deviceToken.destroy()

  getRandomMoinSoundName: ->
    rand = 0
    while ( rand < 1 || rand > 5 || rand == 2 )
      rand = Math.floor @random 1, 6

    "moin" + rand + ".wav"

  random: (min, max) ->
    min + ( Math.random() * (max - min ) )

  send: (sender, receipient, callback) ->
    if !@isRunning
      return callback? new Error 'APN service is not running.'
    if !sender?.getPublicModel?
      return callback? new Error 'Must provide database object for "sender".'
    if !receipient?.getPublicModel?
      return callback? new Error 'Must provide database object for "receipient".'

    receipient.getAPNDeviceTokens().complete (err, deviceTokens) =>
      return callback err if !!err

      for token, i in deviceTokens
        tokenBuffer = new Buffer token.uid, 'hex'
        device = new apn.Device tokenBuffer

        push = new apn.Notification()

        # expiry crashes with "value out of bounds"
        #push.expiry = Date.now() / 1000 + 3600 # 1 hr lifetime
        push.badge = 1
        push.alert = "by " + sender.username
        push.sound = @getRandomMoinSoundName()
        push.payload = {
          "sender": sender.getPublicModel()
        }

        @apnConnection.pushNotification push, device

      callback? null, null

module.exports.APNPush = APNPush
