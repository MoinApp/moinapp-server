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
  match = HEROKU_URL.match /postgres:\/\/([^:]+):([^@]+)@([^:]+):(\d+)\/(.+)/
  
  sequelize = new Sequelize match[5], match[1], match[2], {
    dialect:  'postgres',
    protocol: 'postgres',
    port:     match[4],
    host:     match[3],
    logging:  false
  }


User = require('./models/user') sequelize
Session = require('./models/session') sequelize
gcmID = require('./models/gcmID') sequelize
# Relations
# Push IDs
User.hasMany gcmID
gcmID.belongsTo User
# Login sessions
User.hasMany Session
Session.belongsTo User

if !isHeroku()
  sequelize.sync({ force:true }).success () ->
    crypt = require './crypt'
    
    # create dummy user
    User.createUserAndEncryptPassword({
      username: 'sgade',
      password: 'test',
      email: 'a.b@c.d'
    }).complete (err, user) ->
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
  gcmID: gcmID
}
