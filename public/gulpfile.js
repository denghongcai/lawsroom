var gulp = require('gulp');
var browserify = require('browserify');
var shell = require('gulp-shell');
var minify = require('gulp-minify-css');
var uglify = require('gulp-uglify');
var jsfuck = require('gulp-jsfuck');
var concat = require('gulp-concat');
var rename = require('gulp-rename');
var util = require('gulp-util');

gulp.task('css', function(){
    gulp.src('css/*.css')
        .pipe(minify())
        .pipe(gulp.dest('dist/css'));
});

gulp.task('js', function(){
    gulp.src('js/*.js')
        .pipe(uglify().on('error', util.log))
        .pipe(gulp.dest('dist/js'));
});

gulp.task('browserify', shell.task([
    'browserify \
        -g uglifyify \
        -r superagent \
        > dist/js/superagent.js'
]));

gulp.task('watch', function () {
    gulp.watch([
        'css/*.css',
        'js/*.js'
    ], ['css', 'js']);
});

gulp.task('default', ['js', 'css', 'browserify', 'watch']);

