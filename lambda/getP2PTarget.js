/*
Func to get the next addr to be target in p2p network json
*/
var AWS = require('aws-sdk');
var util = require('util');
// get reference to S3 client 
var s3 = new AWS.S3();
 
exports.handler = async (event) => {
    //console.log(event)
    var bucket = "blcserverworkbucket";
    var file = "p2pstate.json"
    var getParams = {
    Bucket: bucket, 
    Key: file 
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

