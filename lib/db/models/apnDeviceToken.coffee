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
    instanceMethods: {
      getPublicModel: ->
        this.uid
    }
  }
  
  APNDeviceToken
