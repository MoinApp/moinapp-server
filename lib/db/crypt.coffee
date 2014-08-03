crypto = require 'crypto'

exports._getHash = (hash, text) ->
  hash = crypto.createHash hash
  
  hash.update text
  
  hash.digest 'hex'

exports.getMD5 = (text) ->
  exports._getHash 'md5', text

exports.getSHA256 = (text) ->
  exports._getHash 'sha256', text
