/*
Need set role:
Para crear un rol de ejecución
Abra la página Roles en la consola de IAM.
Elija Create role.
Cree un rol con las propiedades siguientes.
Trusted entity (Entidad de confianza) – AWS Lambda.
Permisos – AWSLambdaExecute.
Role name (Nombre de rol): lambda-s3-role.
*/
// [{"Ip":"","Port":"","Key":"","PrevKey":"key1","Configure":false}]
// dependencies
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
      Bucket: bucket, // your bucket name,
      Key: file // path to the object you're looking for
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
            obj.push({"Ip":newIp,"Port":newPort,"Key":newKey,"PrevKey":"","Configure":true});
          }
          if(needSet == true){
            obj.push({"Ip":"","Port":"","Key":"","PrevKey":newKey,"Configure":false});
            //configuring parameters
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

