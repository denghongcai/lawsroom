var pages = document.querySelector('#pages');
var room = document.querySelector('#x-room');
var tb = document.querySelector('paper-toolbar');

(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','https://www.google-analytics.com/analytics.js','ga');
ga('create', 'UA-77171491-1', 'auto');

function _a(ctx){
    ga('set', 'page', ctx.path);
    ga('send', 'pageview');
}

page('/', function(ctx, next){
    tb.hidden = true;
    pages.select('x-door');
    next();
}, _a);
page('/random', function(ctx, next){
    tb.hidden = false;
    pages.select('random-room');
    next();
}, _a);
page('/room/:id', function(ctx, next){
    tb.hidden = false;
    room.roomId = ctx.params.id;
    pages.select('x-room');
    next();
}, _a);
page('*', function(ctx, next){
    tb.hidden = true;
    pages.select('x-door');
    next();
}, _a);

page({
    hashbang: false
});

