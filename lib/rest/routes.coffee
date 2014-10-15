# Module to bundle all routes in one object

module.exports = {
  index: require('./routes/index'),
  moin: require('./routes/moin'),
  session: require('./routes/session'),
  user: require('./routes/user')
}
