var AWS = require('aws-sdk');
var s3 = new AWS.S3();
 
exports.handler = async (event) => {
    var obj = JSON.parse(event);
    if(obj != null && obj != ""){
      var splitnode = obj.newnode.split("/");
      var ip = splitnode[0];
      var port = splitnode[1];
      var newKey = splitnode[2];
      var newPrevkey = splitnode[3];
      var bucket = "blcserverworkbucket";
      var file = "p2pstate.json";
      var getParams = {
      Bucket: bucket, // your bucket name,
      Key: file // path to the object you're looking for
      };
      return await s3.getObject(getParams).promise()
      .then((res) => {
          var obj = JSON.parse(res.Body); 
          var configure = "ko";
          if(obj.length > 0){
            for (let i = 0; i < obj.length; i++){
              if(obj[i].Configure == true && obj[i].Ip == ip && obj[i].Port == port ){
                  configure = "ok";
              }
            }
          }
          const response = {
            statusCode: 200,
            body: configure,
          };
          return response;
      })
      .catch((err) => {
        const response = {
          statusCode: 202,
          body: "ko"+err,
        };
        return response;
      });
  } else {
    const response = {
      statusCode: 202,
      body: "ko",
    };
    return response;
  }
};

