var gulp = require('gulp');
var browserify = require('browserify');
var shell = require('gulp-shell');
var css = require('gulp-minify-css');
var js = require('gulp-uglify');
var concat = require('gulp-concat');
var rename = require('gulp-rename');

gulp.task('css', function(){
    gulp.src('css/*.css')
        .pipe(css())
        .pipe(gulp.dest('dist/css'));
});

gulp.task('js', function(){
    gulp.src('js/*.js')
        .pipe(js())
        .pipe(gulp.dest('dist/js'));
});

gulp.task('bundle', shell.task([
    'browserify \
        -g uglifyify \
        -r keyboardjs \
        > dist/bundle.js'
]));

gulp.task('watch', function () {
    gulp.watch([
        'css/*.css',
        'js/*.js'
    ], ['css', 'js']);
});

gulp.task('default', ['js', 'css', 'bundle', 'watch']);

