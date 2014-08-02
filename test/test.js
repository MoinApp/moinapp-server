var http = require('http');

function postTo(path, callback) {
  return http.request({
    hostname: 'localhost',
    port: 3000,
    path: path,
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    }
  }, callback);
}

var req = postTo('/user/session', function(res) {
  res.on('data', function(chunk) {
    var data = JSON.parse(chunk.toString());
    
    sessionToken = data.session;
    console.log("SessionToken:", sessionToken);
    
    http.get('http://localhost:3000/user/sgade?session=' + sessionToken, function(res) {
      res.on('data', function(chunk) {
        var userInfo = JSON.parse(chunk.toString());
        
        console.log("sgade:", userInfo);
        
        var req = postTo('/moin?session=' + sessionToken, function(res) {
          res.on('data', function(chunk) {
            console.log("/moin", chunk.toString());
          });
        });
        var data = {
          to: userInfo.id
        };
        req.write(JSON.stringify(data));
        req.end();
        
      });
    });
  });
});

var data = {
  username: 'sgade',
  password: 'test'
};
req.write(JSON.stringify(data));
req.end();
