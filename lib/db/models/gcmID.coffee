Sequelize = require 'sequelize'

module.exports = (sequelize) ->
  
  gcmID = sequelize.define 'gcmID', {
    uid: {
      type: Sequelize.STRING,
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
  
  gcmID
