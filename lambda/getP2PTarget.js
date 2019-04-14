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
var util = require('util');
// get reference to S3 client 
var s3 = new AWS.S3();
 
exports.handler = async (event) => {
    //console.log(event)
    var bucket = "blcserverworkbucket";
    var file = "p2pstate.json"
    var getParams = {
    Bucket: bucket, // your bucket name,
    Key: file // path to the object you're looking for
    }
    return await s3.getObject(getParams).promise()
    .then((res) => {
        var obj = JSON.parse(res.Body); 
        var strNode = "empty";
        var prevKey = "";
        for (let i = 0; i < obj.length; i++){
            if(obj[i].Configure == false && obj[i].Ip == "" && obj[i].Port == ""){
                prevKey = obj[i].PrevKey;
            }
        }
        for (let i = 0; i < obj.length; i++){
            if(obj[i].Key == prevKey){
                strNode = obj[i].Ip +  ":" + obj[i].Port + ":" + obj[i].Key;
            }
        }
        const response = {
            statusCode: 200,
            body: strNode,
        };
        return response;
    })
    .catch((err) => {
        return err;
    });
};

