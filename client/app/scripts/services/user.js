'use strict';

/**
 * @ngdoc service
 * @name clientApp.User
 * @description
 * # User
 * Factory in the clientApp.
 */
angular.module('clientApp')
    .factory('User', ['$resource', function ($resource) {
        return $resource('/api/users/:id', {
            id: '@id'
        }, {
            update: {
                method: 'PUT'
            }
        });
    }]);