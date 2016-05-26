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
    if(bowser.firefox){
        return true;
    }
    return false;
}
if(!support()){
    document.body.style.backgroundImage = "none";
    document.body.innerHTML = '<div style="text-align:center;">请使用Chrome/Firefox/Opera浏览器</div>';
};

