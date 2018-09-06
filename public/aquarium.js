var app = angular.module('aquarium', [])

app.controller('state', function ($scope, $http, $interval) {
	$scope.getState = function () {
		$http.get('api/state').
		then(function (response) {
			$scope.aq_state = response.data;
			$scope.aq_error = null
		}, function (response) {
			$scope.aq_state = null
			$scope.aq_error = response.status;
		});
	}

	$scope.getState()
	getPeriodically = $interval($scope.getState, 30000);

	$scope.setSensorName = function (id, name) {
		$http.put('api/sensor/' + id + '/name', name).
			then(function (response) {
			});
	};
	$scope.setRelayName = function (id, name) {
		$http.put('api/relay/' + id + '/name', name).
			then(function (response) {
			});
	};
	$scope.setActionName = function (id, name) {
		$http.put('api/action/' + id + '/name', name).
			then(function (response) {
			});
	};
});