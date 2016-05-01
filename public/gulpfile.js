var gulp = require('gulp');
var browserify = require('browserify');
var shell = require('gulp-shell');
var css = require('gulp-minify-css');
var uglify = require('gulp-uglify');
var jsfuck = require('gulp-jsfuck');
var concat = require('gulp-concat');
var rename = require('gulp-rename');
var util = require('gulp-util');

gulp.task('css', function(){
    gulp.src('css/*.css')
        .pipe(css())
        .pipe(gulp.dest('dist/css'));
});

gulp.task('js', function(){
    gulp.src('js/*.js')
        .pipe(uglify().on('error', util.log))
        .pipe(gulp.dest('dist/js'));
});

//gulp.task('keyboardjs', shell.task([
    //'browserify \
        //-g uglifyify \
        //-r keyboardjs \
        //> dist/js/keyboardjs.js'
//]));

gulp.task('watch', function () {
    gulp.watch([
        'css/*.css',
        'js/*.js'
    ], ['css', 'js']);
});

gulp.task('default', ['js', 'css', 'watch']);

