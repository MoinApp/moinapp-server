Sequelize = require 'sequelize'

HEROKU_URL = process.env.HEROKU_POSTGRESQL_JADE_URL
isHeroku = ->
  return HEROKU_URL?


if !isHeroku()
  # local
  dbConfig = {
    name: 'db.sqlite',
    username: 'admin',
    password: 'admin'
  }

  sequelize = new Sequelize dbConfig.name, dbConfig.username, dbConfig.password, {
    dialect: 'sqlite',
    storage: dbConfig.name,
    #logging: console.log
    logging: false
  }
else
  sequelize = new Sequelize HEROKU_URL, {
    dialect:  'postgres',
    protocol: 'postgres',
    dialectOptions: {
        ssl: true
    }
  }


User = require('./models/user') sequelize
Session = require('./models/session') sequelize
gcmID = require('./models/gcmID') sequelize
APNDeviceToken = require('./models/apnDeviceToken') sequelize
# Relations
# Push IDs
User.hasMany gcmID
gcmID.belongsTo User
User.hasMany APNDeviceToken
APNDeviceToken.belongsTo User
# User's recents
User.hasMany User, { as: 'Recents', joinTableName: 'userRecents' }
# Login sessions
User.hasMany Session
Session.belongsTo User

if !isHeroku()
  sequelize.sync({ force:false }).success () ->
    crypt = require './crypt'

    # create dummy user
    User.createUserAndEncryptPassword({
      username: 'sgade',
      password: 'testtest',
      email: 'a.b@c.d'
    }).complete (err, user) ->
      if user
        gcmID.create({
          uid: 'blub'
        }).complete (err, gcmId) ->
          user.addGcmID gcmId
else
  sequelize.sync().complete (err) ->
    if !!err
      console.log "Error syncing database:", err

module.exports = {
  Sequelize: Sequelize,
  sequelize: sequelize,

  User: User,
  Session: Session,
  gcmID: gcmID,
  APNDeviceToken: APNDeviceToken
}
