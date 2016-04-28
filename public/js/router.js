var pages = document.querySelector('#pages');

page('/', function(ctx, next){
    pages.select('random-chat');
});
page('/room/:id', function(ctx, next){
    document.querySelector('x-room').roomId = ctx.params.id;
    pages.select('x-room');
});

page('*', function(ctx, next){
    pages.select('random-chat');
});
page({
    hashbang: true
});

