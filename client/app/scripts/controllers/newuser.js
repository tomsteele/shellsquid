'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:NewuserCtrl
 * @description
 * # NewuserCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('NewUserCtrl', function ($location, User) {
        var vm = this;
        vm.email = '';
        vm.password = '';
        vm.addUser = function () {
            var user = new User();
            user.email = vm.email;
            user.password = vm.password;
            User.save(user, function () {
                $location.path('/users');
            });
        };

    });