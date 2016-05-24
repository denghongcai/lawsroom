window.logs = [];
function log(s){
    if(!window._debug){
        return;
    }
    console.log(s);
    logs.push(s);
}
