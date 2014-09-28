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
    emailHash: Sequelize.STRING, # md5
    
    session: Sequelize.STRING # uuid v1
  }, {
    classMethods: {
      createUser: (properties) ->
        properties.uid = uuid.v4()
        
        User.create(properties)
    },
    instanceMethods: {
      getPublicModel: ->
        {
          id: this.uid,
          username: this.username,
          email_hash: this.emailHash
        }
    }
  }
  
  User
