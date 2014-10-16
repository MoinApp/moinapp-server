Sequelize = require 'sequelize'
uuid = require 'node-uuid'

module.exports = (sequelize) ->
  
  Session = sequelize.define 'Session', {
    uid: {
      type: Sequelize.STRING, # uuid v4
      unique: true
    },
    application: Sequelize.STRING # app identificer
  }, {
    classMethods: {
      createNew: (appID) ->
        properties = {
          uid: uuid.v1() # time-based,
          application: appID
        }
        
        Session.findOrCreate properties
    },
    instanceMethods: {
      getPublicModel: ->
        @uid
    }
  }
  
  Session
