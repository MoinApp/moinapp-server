Sequelize = require 'sequelize'
uuid = require 'node-uuid'
crypt = require '../crypt'

module.exports = (sequelize) ->
  
  User = sequelize.define 'User', {
    uid: {
      type: Sequelize.STRING, # uuid v4
      unique: true
    },
    username: {
      type: Sequelize.STRING,
      unique: true
    },
    password: Sequelize.STRING, # sha256
    emailHash: Sequelize.STRING # md5
  }, {
    classMethods: {
      createUserAndEncryptPassword: (properties) ->
        properties.uid = uuid.v4()
        if !!properties.password
          properties.password = crypt.getSHA256 properties.password
        if !!properties.email
          properties.emailHash = crypt.getMD5 properties.email
          delete properties.email
        
        User.create(properties)
    },
    instanceMethods: {
      getPublicModel: ->
        {
          id: this.uid,
          username: this.username,
          email_hash: this.emailHash
        }
      isValidPassword: (password) ->
        @password == crypt.getSHA256 password
    }
  }
  
  User
