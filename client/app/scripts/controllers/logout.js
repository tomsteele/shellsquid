'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:LogoutCtrl
 * @description
 * # LogoutCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('LogoutCtrl', function ($location, AuthService) {
        AuthService.logout();
        $location.path('/login');
    });