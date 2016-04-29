var pages = document.querySelector('#pages');
var room = document.querySelector('#x-room');

page('/', function(ctx, next){
    pages.select('x-door');
});
page('/random', function(ctx, next){
    pages.select('random-room');
});
page('/room/:id', function(ctx, next){
    room.roomId = ctx.params.id;
    room.i = Date.now().toString();
    pages.select('x-room');
});

page('*', function(ctx, next){
    pages.select('x-door');
});
page({
    hashbang: true
});

