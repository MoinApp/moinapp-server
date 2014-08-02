Sequelize = require 'sequelize'
uuid = require 'node-uuid'

dbConfig = {
  name: 'db.sqlite',
  username: 'admin',
  password: 'admin'
}

sequelize = new Sequelize dbConfig.name, dbConfig.username, dbConfig.password, {
  dialect: 'sqlite',
  storage: dbConfig.name,
  logging: console.log # for now
}

User = sequelize.define 'User', {
  uid: Sequelize.STRING,
  username: Sequelize.STRING,
  password: Sequelize.STRING,
  
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
        username: this.username
      }
  }
}

gcmID = sequelize.define 'gcmID', {
  uid: Sequelize.STRING
}

User.hasMany gcmID

sequelize.sync({force: true}).success () ->
  User.createUser({
    username: 'sgade',
    password: 'test'
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
