Sequelize = require 'sequelize'

module.exports = (sequelize) ->
  
  APNDeviceToken = sequelize.define 'APNDeviceToken', {
    uid: {
      type: Sequelize.STRING,
      unique: true,
      validate: {
        notEmpty: true
      }
    }
  }, {
    classMethods: {
      createNew: (apnToken) ->
        properties = {
          uid: apnToken
        }

        APNDeviceToken.findOrCreate properties
    }
    instanceMethods: {
      getPublicModel: ->
        this.uid
    }
  }
  
  APNDeviceToken
