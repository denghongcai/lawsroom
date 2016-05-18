function support(){
    if(/micromessenger/.test(navigator.userAgent.toLowerCase())){
        return false;
    }
    if(bowser.blink){
        return true;
    }
    if(bowser.chrome){
        return true;
    }
    return false;
}
if(!support()){
    body.style.backgroundImage = "none";
    document.body.innerHTML = '<div style="text-align:center;">请使用Chrome浏览器</div>';
};

