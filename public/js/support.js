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
    document.body.innerHTML = '<div style="text-align:center;">Please love Chrome/Firefox/Opera</div>';
};

