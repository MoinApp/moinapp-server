Sequelize = require 'sequelize'

module.exports = (sequelize) ->
  
  gcmID = sequelize.define 'gcmID', {
    uid: Sequelize.STRING
  }, {
    instanceMethods: {
      getPublicModel: ->
        this.uid
    }
  }
  
  gcmID
