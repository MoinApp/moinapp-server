Sequelize = require 'sequelize'

dbConfig = {
  name: 'db.sqlite',
  username: 'admin',
  password: 'admin'
}

sequelize = new Sequelize dbConfig.name, dbConfig.username, dbConfig.password, {
  dialect: 'sqlite',
  storage: dbConfig.name,
  logging: true # for now
}

User = sequelize.define 'User', {
  username: Sequelize.STRING,
  password: Sequelize.STRING,
  gcmIDS:
}

GCM_ID = sequelize.define 'GCM_ID', {
  uid: Sequelize.STRING
}

User.hasMany GCM_ID
GCM_ID.hasOne User

module.exports = {
  Sequelize: Sequelize,
  sequelize: sequelize,
  
  User: User,
  GCM_ID: GCM_ID
}
