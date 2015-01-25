'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:HelpCtrl
 * @description
 * # HelpCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('HelpCtrl', function ($scope) {
        $scope.awesomeThings = [
            'HTML5 Boilerplate',
            'AngularJS',
            'Karma'
        ];
    });