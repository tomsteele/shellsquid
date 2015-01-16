'use strict';

/**
 * @ngdoc service
 * @name clientApp.record
 * @description
 * # record
 * Factory in the clientApp.
 */
angular.module('clientApp')
    .factory('Record', function ($resource) {
        return $resource('/api/records/:id', {
            id: '@id'
        }, {
            update: {
                method: 'PUT'
            }
        });
    });