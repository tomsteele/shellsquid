'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:LoginCtrl
 * @description
 * # LoginCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('LoginCtrl', function ($location, AuthService) {
        var vm = this;
        vm.email = '';
        vm.password = '';
        vm.error = '';
        vm.login = function () {
            AuthService.login({
                    email: vm.email,
                    password: vm.password
                })
                .success(function (data) {
                    localStorage.setItem('token', data.token);
                    return $location.path('/records');
                })
                .error(function () {
                    vm.password = '';
                    vm.error = 'Invalid username or password';
                });
        };
        vm.dismiss = function () {
            vm.error = '';
        };
    });