'use strict';

/**
 * @ngdoc overview
 * @name clientApp
 * @description
 * # clientApp
 *
 * Main module of the application.
 */
angular
  .module('clientApp', [
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'angular-jwt'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/help', {
        templateUrl: 'views/help.html',
        controller: 'HelpCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/records', {
        templateUrl: 'views/records.html',
        controller: 'RecordsCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/records/new', {
        templateUrl: 'views/new_record.html',
        controller: 'NewRecordCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/records/:id', {
        templateUrl: 'views/record.html',
        controller: 'RecordCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/users', {
        templateUrl: 'views/users.html',
        controller: 'UsersCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/users/new', {
        templateUrl: 'views/new_user.html',
        controller: 'NewUserCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/users/:id', {
        templateUrl: 'views/user.html',
        controller: 'UserCtrl',
        controllerAs: 'vm',
        data: {
          auth: true
        }
      })
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'LoginCtrl',
        controllerAs: 'vm',
        data: {
          auth: false
        }
      })
      .when('/logout', {
        template: ' ',
        controller: 'LogoutCtrl',
        data: {
          auth: false
        }
      })
      .otherwise({
        redirectTo: '/records'
      });
  })
  .config(function ($httpProvider, jwtInterceptorProvider) {

    jwtInterceptorProvider.tokenGetter = function () {
      return localStorage.getItem('token');
    };

    $httpProvider.interceptors.push('jwtInterceptor');
  })
  .run(function ($rootScope, $location, AuthService) {
    $rootScope.$on('$routeChangeStart', function (evt, next) {
      if (next.data && next.data.auth) {
        if (!AuthService.isAuthenticated()) {
          evt.preventDefault();
          $location.path('/login');
        }
      }
    });
  });