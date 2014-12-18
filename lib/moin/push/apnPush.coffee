apn = require 'apn'
db = require '../../db/'

class APNPush
  constructor: (moinController) ->
    # do some init here

    @apnConnection = new apn.Connection {}

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

  send: (sender, receipient, callback) ->
    if sender?.getPublicModel?
      return callback new Error 'Must provide database object for "sender".'
    if receipient?.getPublicModel?
      return callback new Error 'Must provide database object for "receipient".'

    # TODO: Push

    iOSDeviceTokens = [] # recipient.getIOsDeviceTokens()
    for token, i in iOSDeviceTokens
      device = new apn.Device token

      push = new apn.Notification()

      push.expiry = Date.now() / 1000 + 3600 # 1 hr lifetime
      push.badge = 1
      push.alert = "Moin"
      push.payload = {
        "sender": sender.getPublicModel()
      }

      @apnConnection.pushNotification push, device

    callback new Error 'Not implemented.', null

module.exports.APNPush = APNPush
