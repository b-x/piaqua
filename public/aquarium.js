var app = angular.module('aquarium', [])

app.controller('state', function ($scope, $http) {
	$http.get('api/state').
		then(function (response) {
			$scope.aq_state = response.data;
		}, function (response) {
			$scope.aq_error = response.status;
		});

	$scope.setSensorName = function (r) {
		$http.put('api/sensor/' + r.id + '/name', r.name).
			then(function (response) {
			});
	};
	$scope.setRelayName = function (r) {
		$http.put('api/relay/' + r.id + '/name', r.name).
			then(function (response) {
			});
	};
});