/*
Func to set the input addr into p2p network json
*/
var AWS = require('aws-sdk');
var s3 = new AWS.S3();
 
exports.handler = async (event) => {
    var obj = JSON.parse(event);
    if(obj != null && obj != ""){
      var splitnode = obj.newnode.split("/");
      var newIp = splitnode[0];
      var newPort = splitnode[1];
      var newKey = splitnode[2];
      var newPrevkey = splitnode[3];
      var bucket = "blcserverworkbucket";
      var file = "p2pstate.json";
      var getParams = {
      Bucket: bucket, 
      Key: file 
      };
      return await s3.getObject(getParams).promise()
      .then((res) => {
          var obj = JSON.parse(res.Body); 
          var needSet = false;
          if(obj.length > 0){
            for (let i = 0; i < obj.length; i++){
              if(obj[i].Configure == false && obj[i].PrevKey == newPrevkey){
                  obj[i].Key = newKey;
                  obj[i].Ip = newIp;
                  obj[i].Port = newPort;
                  obj[i].Configure = true;
                  needSet = true;
              }
            }
          } else {
            needSet = true;
            obj.push({"Ip":newIp,"Port":newPort,"Key":newKey,"PrevKey":"initial","Configure":true});
          }
          if(needSet == true){
            obj.push({"Ip":"","Port":"","Key":"","PrevKey":newKey,"Configure":false});
            var uploadParams = {
              Bucket: bucket,
              Body : JSON.stringify(obj),
              Key : file
            };
              s3.upload(uploadParams, function (err, data) {
              //handle error
              if (err) {
                return "err";
              }
              //success
              if (data) {
                return "ok";
              }
            });
            const response = {
              statusCode: 201,
              body: "ok"
            };
            return response;
          } else {
            const response = {
              statusCode: 200,
              body: "no need Upload",
            };
            return response;
          }
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

