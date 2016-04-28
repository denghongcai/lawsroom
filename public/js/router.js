var pages = document.querySelector('#pages');

page('/', function(ctx, next){
    pages.select('x-door');
});
page('/random', function(ctx, next){
    pages.select('random-room');
});
page('/room/:id', function(ctx, next){
    //document.querySelector('x-room').roomId = ctx.params.id;
    //document.querySelector('x-room').init = Date.now().toString();
    pages.select('x-room');
});

page('*', function(ctx, next){
    pages.select('x-door');
});
page({
    hashbang: true
});

