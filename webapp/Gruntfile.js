module.exports = function (grunt) {
    require('load-grunt-tasks')(grunt);

    grunt.initConfig({
        clean: ['dist'],
        concat: {
            js : {
                src: [
                    'js/app.js',
                    'js/lib/*.js',
                    'js/ctrl/*.js',
                ],
                dest: 'dist/js/app.js'
            },
            css : {
                src: [
                    'css/*.css'
                ],
                dest: 'dist/css/app.css'
            },
            js_vendors: {
                src: [
                    'node_modules/jquery/dist/jquery.js',
                    'node_modules/bootstrap/dist/js/bootstrap.js',
                    'node_modules/angular/angular.js',
                    'node_modules/ng-file-upload/dist/ng-file-upload-shim.js',
                    'node_modules/ng-file-upload/dist/ng-file-upload.js',
                    'node_modules/angular-sanitize/angular-sanitize.js',
                    'node_modules/angular-route/angular-route.js',
                    'node_modules/angular-ui-bootstrap/dist/ui-bootstrap-tpls.js',
                    'node_modules/angular-markdown-directive/markdown.js',
                    'node_modules/underscore/underscore.js',
                    'node_modules/filesize/lib/filesize.js',
                    'node_modules/showdown/dist/showdown.js',
                    'node_modules/clipboard/dist/clipboard.js'
                ],
                dest: 'dist/js/vendor.js'

            },
            css_vendors: {
                src: [
                    'node_modules/bootstrap/dist/css/bootstrap.css',
                    'node_modules/font-awesome/css/font-awesome.css',
                    'css/water_drop.css'
                ],
                dest: 'dist/css/vendor.css'
            }
        },
        copy: {
            index: {
                files: [{
                    src: [
                        'index.html',
                    ],
                    dest: 'dist/index.html',
                }]
            },
            favicon: {
                files: [{
                    src: [
                        'favicon.ico',
                    ],
                    dest: 'dist/favicon.ico',
                }]
            },
            partials: {
                files: [{
                    expand: true,
                    src: [
                        'partials/*',
                    ],
                    dest: 'dist/partials/',
                    flatten: true
                }]
            },
            images: {
                files: [{
                    expand: true,
                    src: [
                        'img/*',
                    ],
                    dest: 'dist/img/',
                    flatten: true
                }]
            },
            fonts: {
                files: [{
                    expand: true,
                    src: [
                        'node_modules/bootstrap/fonts/*',
                        'node_modules/font-awesome/fonts/*',
                        'fonts/*'
                    ],
                    dest: 'dist/fonts/',
                    flatten: true
                }]
            },
            custom: {
                files: [{
                    src: [
                        'css/custom.css',
                    ],
                    dest: 'dist/css/custom.css',
                },{
                    src: [
                        'js/custom.js',
                    ],
                    dest: 'dist/js/custom.js',
                }]
            },
        },
        ngAnnotate: {
            options: {
                singleQuotes: true
            },
            all: {
                files: {
                    'dist/js/app.js': ['dist/js/app.js'],
                    'dist/js/vendor.js': ['dist/js/vendor.js']
                }
            }
        },
        uglify: {
            options: {
                mangle: false,
                compress: false,
                report: true,
                sourceMap: true
            },
            javascript: {
                files: {
                    'dist/js/app.js': ['dist/js/app.js'],
                    'dist/js/vendor.js': ['dist/js/vendor.js'],
                }
            }

        },
        cssmin: {
            options: {
                keepSpecialComments: 0
            },
            combine: {
                files: {
                    'dist/css/vendor.css': ['dist/css/vendor.css']
                }
            }
        }
    });

    grunt.registerTask('default', ['clean', 'concat', 'copy', 'ngAnnotate', 'uglify', 'cssmin']);
};