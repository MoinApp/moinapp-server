Sequelize = require 'sequelize'
uuid = require 'node-uuid'

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
    logging: console.log
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

gcmID = sequelize.define 'gcmID', {
  uid: Sequelize.STRING
}, {
  instanceMethods: {
    getPublicModel: ->
      this.uid
  }
}

User.hasMany gcmID

sequelize.sync().success () ->
  User.createUser({
    username: 'sgade',
    password: 'test',
    emailHash: '1de9ab0eb9b23a2da38f6857de628625'
  }).complete (err, user) ->
    
    gcmID.create({
      uid: 'blub'
    }).complete (err, gcmId) ->
      user.addGcmID gcmId

module.exports = {
  Sequelize: Sequelize,
  sequelize: sequelize,
  
  User: User,
  gcmID: gcmID
}
