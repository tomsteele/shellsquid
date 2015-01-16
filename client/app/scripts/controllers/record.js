'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:RecordCtrl
 * @description
 * # RecordCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('RecordCtrl', function ($location, $routeParams, Record, User) {
        var vm = this;
        vm.record = {};
        vm.users = [];
        vm.owner = {};

        vm.findOwner = function () {
            for (var i = 0; i < vm.users.length; i++) {
                var user = vm.users[i];
                if (user._id === vm.record.owner._id) {
                    return user;
                }
            }
        };

        vm.updateRecord = function () {
            for (var i = 0; i < vm.users.length; i++) {
                var user = vm.users[i];
                if (user.email === vm.owner.email) {
                    vm.record.owner = user;
                    break;
                }
            }
            vm.record.$update(function () {});
        };
        vm.clearClientsFromRecord = function () {
            vm.record.clients = [];
            vm.record.$update(function () {});
        };
        vm.deleteRecord = function () {
            vm.record.$delete(function () {
                $location.path('/records');
            });
        };
        vm.blacklistRecord = function () {
            vm.record.blacklist = true;
            vm.record.$update(function () {});
        };
        vm.unblacklistRecord = function () {
            vm.record.blacklist = false;
            vm.record.$update(function () {});
        };

        function init() {
            vm.record = Record.get({
                id: $routeParams.id
            }, function () {
                vm.users = User.query({}, function () {
                    vm.owner = vm.findOwner();
                });
            });
        }
        init();
    });