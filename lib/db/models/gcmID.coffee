Sequelize = require 'sequelize'

module.exports = (sequelize) ->

  gcmID = sequelize.define 'gcmID', {
    uid: {
      type: Sequelize.STRING,
      unique: true,
      validate: {
        notEmpty: true
      }
    }
  }, {
    classMethods: {
      createNew: (gcm) ->
        properties = {
          uid: gcm
        }

        gcmID.findOrCreate properties
    }
    instanceMethods: {
      getPublicModel: ->
        this.uid

    }
  }

  gcmID
