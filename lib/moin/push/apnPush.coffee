apn = require 'apn'
db = require '../../db/'

class APNPush
  constructor: (pfxBuffer, moinController) ->
    # do some init here

    # see https://github.com/argon/node-apn/blob/master/doc/connection.markdown
    @apnConnection = new apn.Connection {
      pfx: pfxBuffer,
      connectionTimeout: 5000
    }

    @feedback = new apn.Feedback {
      "interval": ( 6 * 60 * 60 ) # check every 6 hrs
    }
    @feedback.on 'feedback', (time, deviceTokenBuffer) =>
      deviceToken = deviceTokenBuffer.toString()
      @deleteDeviceToken deviceToken

    moinController.on 'moin', (sender, receipient) =>
      @send sender, receipient

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

  send: (sender, receipient, callback) ->
    if sender?.getPublicModel?
      return callback new Error 'Must provide database object for "sender".'
    if receipient?.getPublicModel?
      return callback new Error 'Must provide database object for "receipient".'

    # TODO: Push

    receipient.getAPNDeviceTokens().complete (err, deviceTokens) =>
      return callback err if !!err

      for token, i in deviceTokens
        device = new apn.Device token

        push = new apn.Notification()

        push.expiry = Date.now() / 1000 + 3600 # 1 hr lifetime
        push.badge = 1
        push.alert = "Moin"
        push.payload = {
          "sender": sender.getPublicModel()
        }

        @apnConnection.pushNotification push, device

    callback null, null

module.exports.APNPush = APNPush
