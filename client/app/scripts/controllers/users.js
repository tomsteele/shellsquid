'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:UsersCtrl
 * @description
 * # UsersCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('UsersCtrl', function (User, MessengerService) {
        var vm = this;
        vm.users = [];

        function init() {
            vm.users = User.query({}, function () {}, MessengerService.error);
        }
        init();
    });