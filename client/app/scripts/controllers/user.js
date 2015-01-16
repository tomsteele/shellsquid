'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:UserCtrl
 * @description
 * # UserCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('UserCtrl', function ($http, $routeParams, $location, User, Record) {
        var vm = this;
        vm.user = {};
        vm.password = '';
        vm.confirmPassword = '';
        vm.records = [];
        vm.error = '';
        vm.success = '';

        vm.changePassword = function () {
            if (vm.password !== vm.confirmPassword) {
                vm.error = 'Passwords must match';
                return;
            }
            $http.post('/api/users/' + $routeParams.id, {
                password: vm.password
            }).success(function () {
                vm.success = 'Password updated';
                vm.error = '';
                clean();
            }).error(function () {
                vm.error = 'Error updating password';
                vm.success = '';
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
                    });
                });
            });
        }

        init();
    });