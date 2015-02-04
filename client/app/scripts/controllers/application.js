'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:ApplicationCtrl
 * @description
 * # ApplicationCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('ApplicationCtrl', function ($rootScope, $http, AuthService) {
        var vm = this;
        vm.info = {};
        var firstRun = false;
        vm.isAuthenticated = function () {
            if (AuthService.isAuthenticated()) {
                if (!firstRun) {
                    init();
                    firstRun = true;
                }
                return true;
            }
            return false;
        };
        vm.dismiss = function () {
            $rootScope.successMessage = '';
            $rootScope.errorMessage = '';
        };

        function init() {
            $http.get('/api/info').success(function (data) {
                vm.info = data;
                return true;
            });
        }

    });