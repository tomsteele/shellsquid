'use strict';

/**
 * @ngdoc service
 * @name clientApp.messenger
 * @description
 * # messenger
 * Service in the clientApp.
 */
angular.module('clientApp')
    .service('MessengerService', function ($rootScope, $timeout, $location, AuthService) {
        return {
            error: function (data) {
                if (data.status && data.status === 401) {
                    AuthService.logout();
                    $location.path('/login');
                }

                if (data.data.error && typeof data.data.error === 'string') {
                    $rootScope.errorMessage = data.data.error;
                } else if (Array.isArray(data.data)) {
                    $rootScope.errorMessage = 'missing one or more required fields';
                } else {
                    $rootScope.errorMessage = 'there was an error during the request';
                }
            },
            success: function (msg) {
                $rootScope.successMessage = msg;
                $timeout(function () {
                    $rootScope.successMessage = '';
                }, 2000);
            }
        };
    });