/*
Func handdle log
*/
var AWS = require('aws-sdk');
var s3 = new AWS.S3();

exports.handler = async (event) => {
    var newlog = "";
    var obj = JSON.parse(event);
    if(obj != null && obj != ""){
      newlog = obj.newnode + "";
      var bucket = "blcserverlessbucket";
      var file = "blockchain-serverless.log";
      var getParams = {
      Bucket: bucket, 
      Key: file 
      };
      return await s3.getObject(getParams).promise()
      .then((res) => {
        var obj = "";
          if(typeof(res.Body) !== 'undefined'){
              obj = res.Body.toString('utf-8');
          } else {
            obj = "";
          }
          var needSet = false;
          var allog = "";
          if(obj != ""){
            allog = obj;
            allog = allog.replace(/(\r\n|\n|\r)/gm, "");
            if(allog.indexOf(newlog) == -1){
              var strlog = allog + newlog;
              var splitlogParser = strlog.split("|#FFFFFF");
              var w = 0; 
              var splitlog = Array();
              for(w = 0; w < splitlogParser.length; w++){
                if(splitlogParser[w] != ""){
                  splitlog.push(splitlogParser[w]);
                }
              }
              var i = 0;
              for (i = 0; i < splitlog.length; i++){
                splitlog[i] = splitlog[i] + "|#FFFFFF";
              }
              splitlog.sort(function (a, b){
                var sa = a.split("|");
                var atime = sa[0];
                var aint = Number(atime);
                var sb = b.split("|");
                var btime = sb[0];
                var bint = Number(btime);
                return aint - bint;
              });
              var j = 0;
              allog = "";
              for (j = 0; j < splitlog.length; j++){
                if(j == splitlog.length - 1){
                  allog = allog + splitlog[j];
                } else {
                  allog = allog + splitlog[j] + "\n";
                }
              }
              needSet = true;
            } else {
              needSet = false;
            }
          } else {
            needSet = true;
            allog = newlog;
          }
          if(needSet == true){
            var uploadParams = {
              Bucket: bucket,
              Body : allog,
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

