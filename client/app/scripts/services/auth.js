'use strict';

/**
 * @ngdoc service
 * @name clientApp.auth
 * @description
 * # auth
 * Factory in the clientApp.
 */
angular.module('clientApp')
    .factory('AuthService', function ($http, jwtHelper) {
        return {
            isAuthenticated: function () {
                var token = localStorage.getItem('token');
                if (!token) {
                    return false;
                }
                try {
                    if (jwtHelper.isTokenExpired(token)) {
                        return false;
                    }
                } catch (e) {
                    return false;
                }
                return true;
            },

            login: function (credentials) {
                return $http.post('/api/token', credentials);
            },
            logout: function () {
                return localStorage.setItem('token', null);
            }
        };
    });
