var AWS = require('aws-sdk');
var s3 = new AWS.S3();
var needSet = false;
 
exports.handler = async (event) => {
    var obj = JSON.parse(event);
    if(obj != null && obj != ""){
      var splitnode = obj.newnode.split("/");
      var key = splitnode[2];
      var bucket = "blcserverworkbucket";
      var file = "p2pstate.json";
      var getParams = {
      Bucket: bucket, // your bucket name,
      Key: file // path to the object you're looking for
      };
      return await s3.getObject(getParams).promise()
      .then((res) => {
        var obj = JSON.parse(res.Body); 
        //var needSet = false;
        if(obj.length > 0){
          deleteFromObject(key,obj);
        }
        if(needSet == true){
          setLastNode(obj);
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

function deleteFromObject(key, obj){
  let nextKey = null;
  let size = obj.length;
  for (let i = 0; i < size; i++){
    if(obj[i].Configure == true && obj[i].Key == key){
      nextKey = obj[i].Key; 
      obj.splice(i,1);
      size--;
      i--;
      needSet = true;
    }
  }
  if(nextKey != null){
    deletePrevKey(nextKey,obj);
  }
}

function deletePrevKey(nextKey, obj){
  let nextPrevKey = null;
  let size = obj.length;
  for (let i = 0; i < size; i++){
    if(obj[i].Configure == true && obj[i].PrevKey == nextKey || obj[i].Configure == false && obj[i].PrevKey == nextKey){
      nextPrevKey = obj[i].Key; 
      obj.splice(i,1);
      size--;
      i--;
    } 
  }

  if(nextPrevKey != null && nextPrevKey != ""){
    deletePrevKey(nextPrevKey, obj);
  }
}

function setLastNode(obj){
  let size = obj.length;
  if(size > 0){
    if(obj[size-1].Configure == true){
      let lastKey = obj[size-1].Key;
      obj.push({"Ip":"","Port":"","Key":"","PrevKey":lastKey,"Configure":false});
    }
  }
}

