var pages = document.querySelector('#pages');
var room = document.querySelector('#x-room');
var tb = document.querySelector('paper-toolbar');

page('/', function(ctx, next){
    tb.hidden = true;
    pages.select('x-door');
});
page('/random', function(ctx, next){
    tb.hidden = false;
    pages.select('random-room');
});
page('/room/:id', function(ctx, next){
    tb.hidden = false;
    room.roomId = ctx.params.id;
    pages.select('x-room');
});

page('*', function(ctx, next){
    tb.hidden = true;
    pages.select('x-door');
});
page({
    hashbang: true
});

