var obj = 
[
  {
      "Ip": "127.0.0.1",
      "Port": "10000",
      "Key": "Qmb2UUE8G2kGaFdFVUo2EfaeJfGXrUddV7RrfLLmJ2Jdda",
      "PrevKey": "initial",
      "Configure": true
  },
  {
      "Ip": "127.0.0.1",
      "Port": "10001",
      "Key": "QmeDzSrU6q8AMrABJKNe9h7xKydit5rSu4BLPfWxdTvZtJ",
      "PrevKey": "Qmb2UUE8G2kGaFdFVUo2EfaeJfGXrUddV7RrfLLmJ2Jdda",
      "Configure": true
  },
  {
      "Ip": "127.0.0.1",
      "Port": "10002",
      "Key": "QmV8TzxJQCVxJ3abCteu94WNF6phpGWRZmsw5bjeLy2zPC",
      "PrevKey": "QmeDzSrU6q8AMrABJKNe9h7xKydit5rSu4BLPfWxdTvZtJ",
      "Configure": true
  },
  {
      "Ip": "",
      "Port": "",
      "Key": "",
      "PrevKey": "QmV8TzxJQCVxJ3abCteu94WNF6phpGWRZmsw5bjeLy2zPC",
      "Configure": false
  }
];


var key = "QmeDzSrU6q8AMrABJKNe9h7xKydit5rSu4BLPfWxdTvZtJ";
console.log(obj);
debugger;
deleteFromObject(key,obj);
setLastNode(obj);
console.log(obj);
function deleteFromObject(key, obj){
    let nextKey = null;
    let size = obj.length;
    for (let i = 0; i < size; i++){
      if(obj[i].Configure == true && obj[i].Key == key){
        nextKey = obj[i].Key; 
        obj.splice(i,1);
        size--;
        i--;
        //needSet = true;
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
      if(obj[i].Configure == true && obj[i].PrevKey == nextKey || obj[i].Configure == false && obj[i].PrevKey == nextKey){
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
        let lastKey = obj[size-1].Key
        obj.push({"Ip":"","Port":"","Key":"","PrevKey":lastKey,"Configure":false});
      }
    }
  }
  
  