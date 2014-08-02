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
  password: Sequelize.STRING
}, {
  classMethods: {
    createUser: (username, password) ->
      uid = uuid.v4()
      
      User.create({
        username: username,
        password: password,
        uid: uid
      })
  },
  instanceMethods: {
    getPublicModel: ->
      {
        id: this.uid,
        name: this.username
      }
  }
}

GCM_ID = sequelize.define 'GCM_ID', {
  uid: Sequelize.STRING
}

User.hasMany GCM_ID
GCM_ID.hasOne User

sequelize.sync({force: true}).success () ->
  User.createUser 'sgade', 'test'

module.exports = {
  Sequelize: Sequelize,
  sequelize: sequelize,
  
  User: User,
  GCM_ID: GCM_ID
}
