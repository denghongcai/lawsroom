function support(){
    if(bowser.blink){
        return true;
    }
    if(bowser.chrome){
        return true;
    }
    return false;
}
if(!support()){
    document.body.innerHTML = '<div style="text-align:center;">请使用最新版的Chrome浏览器</div>';
};

