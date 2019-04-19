var obj = [{"Ip":"127.0.0.1","Port":"10000","Key":"Qmf4c5VLKeoW9QDWb8J7XESncxPqKaNUvk4omTSF72dtBb","PrevKey":"","Configure":true},{"Ip":"127.0.0.1","Port":"10001","Key":"QmedVDCf5MJDgC9BH84EcNVkgpo9h12MoTBXVc7e3yWWZ7","PrevKey":"Qmf4c5VLKeoW9QDWb8J7XESncxPqKaNUvk4omTSF72dtBb","Configure":true},{"Ip":"","Port":"","Key":"","PrevKey":"QmedVDCf5MJDgC9BH84EcNVkgpo9h12MoTBXVc7e3yWWZ7","Configure":false}];
var key = "QmedVDCf5MJDgC9BH84EcNVkgpo9h12MoTBXVc7e3yWWZ7";
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
  
  