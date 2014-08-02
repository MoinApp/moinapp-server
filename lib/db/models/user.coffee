Sequelize = require 'sequelize'
uuid = require 'node-uuid'

module.exports = (sequelize) ->
  
  User = sequelize.define 'User', {
    uid: Sequelize.STRING,
    username: Sequelize.STRING,
    password: Sequelize.STRING,
    emailHash: Sequelize.STRING, # md5
    
    session: Sequelize.STRING
  }, {
    classMethods: {
      createUser: (properties) ->
        console.log 'createUser:', properties
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