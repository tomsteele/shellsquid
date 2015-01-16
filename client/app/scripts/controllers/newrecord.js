'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:NewrecordCtrl
 * @description
 * # NewrecordCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
    .controller('NewRecordCtrl', function ($location, Record) {
        var vm = this;
        vm.error = null;
        vm.fqdn = '';
        vm.handler_host = '';
        vm.handler_port = null;
        vm.handler_protocol = '';
        vm.addRecord = function () {
            var record = new Record();
            record.fqdn = vm.fqdn;
            record.handler_host = vm.handler_host;
            record.handler_port = vm.handler_port;
            record.handler_protocol = vm.handler_protocol;
            Record.save(record, function (data) {
                return $location.path('/records/' + data.id);
            });
        };

    });