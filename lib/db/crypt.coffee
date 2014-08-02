crypot = require 'crypto'

exports.getMD5 = (text) ->
  md5 = crypto.createHash 'md5'
  
  md5.update text
  
  md5.digest 'hex'
