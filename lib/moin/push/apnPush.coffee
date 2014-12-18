apn = require 'apn'

class APNPush
  constructor: (moinController) ->
    # do some init here

    @apnConnection = apn.Connection {}

    moinController.on 'moin', (sender, receipient) =>
      @send sender, receipient

  send: (sender, receipient, callback) ->
    if sender?.getPublicModel?
      return callback new Error 'Must provide database object for "sender".'
    if receipient?.getPublicModel?
      return callback new Error 'Must provide database object for "receipient".'

    # TODO: Push

    iOSDeviceTokens = [] # recipient.getIOsDeviceTokens()
    for token, i in iOSDeviceTokens
      device = apn.Device token

      push = apn.Notification()

      push.expiry = Date.now() / 1000 + 3600 # 1 hr lifetime
      push.badge = 1
      push.alert = "Moin"
      push.payload = {
        "sender": sender.getPublicModel()
      }

      @apnConnection.pushNotification push, device

    callback new Error 'Not implemented.', null

module.exports.APNPush = APNPush
