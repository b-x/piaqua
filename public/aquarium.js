var app = angular.module('aquarium', [])

app.controller('state', function ($scope, $http) {
	$http.get('api/state').
		then(function (response) {
			$scope.aq_state = response.data;
		}, function (response) {
			$scope.aq_error = response.status;
		});

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