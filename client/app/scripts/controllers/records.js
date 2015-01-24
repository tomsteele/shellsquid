'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:RecordsCtrl
 * @description
 * # RecordsCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('RecordsCtrl', function ($scope, Record, MessengerService) {
        var vm = this;
        vm.records = [];

        function init() {
            vm.records = Record.query(function () {}, MessengerService.error);
        }

        init();
    });