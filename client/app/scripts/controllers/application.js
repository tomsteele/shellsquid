'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:ApplicationCtrl
 * @description
 * # ApplicationCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('ApplicationCtrl', function ($rootScope, AuthService) {
        var vm = this;
        vm.isAuthenticated = AuthService.isAuthenticated;
        vm.dismiss = function () {
            $rootScope.successMessage = '';
            $rootScope.errorMessage = '';
        };
    });