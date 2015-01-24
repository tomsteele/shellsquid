'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:UserCtrl
 * @description
 * # UserCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('UserCtrl', function ($http, $routeParams, $location, User, Record, MessengerService) {
        var vm = this;
        vm.user = {};
        vm.password = '';
        vm.confirmPassword = '';
        vm.records = [];

        vm.changePassword = function () {
            if (vm.password !== vm.confirmPassword) {
                vm.error = 'Passwords must match';
                return;
            }
            $http.post('/api/users/' + $routeParams.id, {
                password: vm.password
            }).success(function () {
                MessengerService.success('password updated');
                clean();
            }).error(function (data) {
                MessengerService.error(data);
                clean();
            });

            function clean() {
                vm.password = '';
                vm.confirmPassword = '';
            }
        };

        vm.deleteUser = function () {

            vm.user.$delete(function () {
                $location.path('/users');
            });
        };

        function init() {
            vm.user = User.get({
                id: $routeParams.id
            }, function () {
                var records = Record.query({}, function () {
                    vm.records = records.filter(function (record) {
                        if (record.owner.id === vm.user.id) {
                            return record;
                        }
                    }, MessengerService.error);
                });
            });
        }

        init();
    });