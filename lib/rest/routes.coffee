# Module to bundle all routes in one object

module.exports = {
  v200: {
    version: "2.0.0",
    
    index: require('./api-v2/index'),
    moin: require('./api-v2/moin'),
    session: require('./api-v2/session'),
    user: require('./api-v2/user')
  }
}
