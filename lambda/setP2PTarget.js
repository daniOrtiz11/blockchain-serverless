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

// dependencies
var AWS = require('aws-sdk');
var s3 = new AWS.S3();
 
exports.handler = async (event) => {
    var newnode = event.newnode;
    if(newnode != null && newnode != ""){
        var splitnode = newnode.split("/");
    }
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
        for (let i = 0; i < obj.length; i++){
            if(obj[i].Configure == false && obj[i].PrevKey == newPrevkey){
                obj[i].Key = newKey;
                obj[i].Ip = newIp;
                obj[i].Port = newPort;
                obj[i].Configure = true;
            }
        }
        obj.push({"Ip":"","Port":"","Key":"","PrevKey":newKey,"Configure":false});
        //configuring parameters
        var uploadParams = {
          Bucket: bucket,
          Body : JSON.stringify(obj),
          Key : file
        };
        var ok = "";
        s3.upload(uploadParams, function (err, data) {
          //handle error
          if (err) {
            return err
          }
          //success
          if (data) {
            return "ok";
          }
        });
        const response = {
            statusCode: 200,
            body: "ok",
        };
        return response;
    })
    .catch((err) => {
        return err;
    });
};

