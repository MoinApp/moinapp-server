apn = require 'apn'

class APNPush
  constructor: (moinController) ->
    # do some init here

    moinController.on 'moin', (sender, receipient) =>
      @send sender, receipient

  send: (sender, receipient, callback) ->
    if sender?.getPublicModel?
      return callback new Error 'Must provide database object for "sender".'
    if receipient?.getPublicModel?
      return callback new Error 'Must provide database object for "receipient".'

    # TODO: Push

    callback new Error 'Not implemented.', null

module.exports.APNPush = APNPush
