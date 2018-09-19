var app = angular.module('aquarium', [])

app.directive('aqDuration', function () {
	return {
		restrict: 'A',
		require: 'ngModel',
		link: function (scope, element, attrs, ngModel) {
			if (ngModel) {
				ngModel.$parsers.push(function (value) {
					return value.valueOf() * 1000000;
				});
				ngModel.$formatters.push(function (value) {
					return new Date(value / 1000000);
				});
				ngModel.$overrideModelOptions({
					timezone: 'UTC',
					timeSecondsFormat: 'ss',
				});
			}
		}
	};
});

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

	$scope.toggleAction = function (id) {
		$http.put('api/action/' + id + '/toggle').
			then(function (response) {
			})
	}

	$scope.updateAction = function (id, action) {
		$http.put('api/action/' + id, action).
			then(function (response) {
				$scope.showState()
			})
	}

	$scope.addAction = function (action) {
		$http.post('api/action', action).
			then(function (response) {
				$scope.showState()
			})
	}

	$scope.removeAction = function (id) {
		$http.delete('api/action/' + id).
			then(function (response) {
				$scope.showState()
			})
	}

	$scope.showState = function () {
		$scope.getState()
		$scope.aq_form = 'state'
	}

	$scope.editAction = function (id) {
		$scope.aq_edit_action_id = id
		if (id >= 0) {
			$scope.aq_edit_action = angular.copy($scope.aq_state.actions[id])
		} else {
			$scope.aq_edit_action = { duration: 0 }
		}
		$scope.aq_edit_action_relays = $scope.aq_state.relays.map(x => x.name)
		$scope.aq_edit_action_buttons = Array.from(new Array($scope.aq_state.buttons), (x, i) => i + 1)
		$scope.aq_form = 'action'
	}

	$scope.editRelayTasks = function (id) {
		$scope.aq_form = 'tasks'
	}

	$scope.showState()
	getPeriodically = $interval($scope.getState, 30000);
});