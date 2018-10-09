const weekDays = ['Pn', 'Wt', 'Śr', 'Cz', 'Pt', 'Sb', 'Nd'];

function weekDayBit(day) {
  const weekDaysOrder = [1, 2, 3, 4, 5, 6, 0];
  return 1 << weekDaysOrder[day];
}

const app = angular.module('aquarium', []);

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

app.filter('aqDuration', function ($filter) {
  return function (value) {
    return $filter('date')(value / 1000000, 'HH:mm:ss', 'UTC');
  };
});

app.filter('aqWeekdays', function () {
  return function (value) {
    if (value == 127) {
      return "codziennie";
    }
    return weekDays
      .filter((d, i) => value & weekDayBit(i))
      .map(d => d.toLowerCase())
      .join(', ');
  };
});

app.controller('state', function ($scope, $http, $interval) {
  $scope.getState = function () {
    $http.get('api/state').
      then(function (response) {
        $scope.aq_state = response.data;
        $scope.aq_error = null;
      }, function (response) {
        $scope.aq_state = null;
        $scope.aq_error = response.status;
      });
  };

  $scope.setSensorName = function (id, name) {
    $http.put('api/sensor/' + id + '/name', name).
      then(function (response) {
        $scope.showState();
      });
  };

  $scope.setRelayName = function (id, name) {
    $http.put('api/relay/' + id + '/name', name).
      then(function (response) {
        $scope.showState();
      });
  };

  $scope.toggleAction = function (id) {
    $http.put('api/action/' + id + '/toggle').
      then(function (response) {
        $scope.getState();
      });
  };

  $scope.updateAction = function (id, action) {
    $http.put('api/action/' + id, action).
      then(function (response) {
        $scope.showState();
      });
  };

  $scope.addAction = function (action) {
    $http.post('api/action', action).
      then(function (response) {
        $scope.showState();
      });
  };

  $scope.removeAction = function (id) {
    $http.delete('api/action/' + id).
      then(function (response) {
        $scope.showState();
      });
  };

  $scope.updateTask = function (relay_id, id, task) {
    $http.put('api/relay/' + relay_id + '/task/' + id, task).
      then(function (response) {
        $scope.showRelay();
      });
  };

  $scope.addTask = function (relay_id, task) {
    $http.post('api/relay/' + relay_id + '/task', task).
      then(function (response) {
        $scope.showRelay();
      });
  };

  $scope.removeTask = function (relay_id, id) {
    $http.delete('api/relay/' + relay_id + '/task/' + id).
      then(function (response) {
        $scope.showRelay();
      });
  };

  $scope.showState = function () {
    $scope.getState();
    $scope.aq_form = 'state';
  };

  $scope.showRelay = function () {
    $scope.getState();
    $scope.aq_form = 'relay';
  };

  $scope.editAction = function (id) {
    $scope.aq_edit_action_id = id;
    if (id >= 0) {
      $scope.aq_edit_action = angular.copy($scope.aq_state.actions[id]);
    } else {
      $scope.aq_edit_action = { duration: 0 };
    }
    $scope.aq_edit_action_relays = $scope.aq_state.relays.map(x => x.name);
    $scope.aq_edit_action_buttons = Array.from(new Array($scope.aq_state.buttons), (x, i) => i + 1);
    $scope.aq_form = 'action';
  };

  $scope.editSensor = function (id) {
    $scope.aq_edit_sensor_id = id;
    $scope.aq_edit_sensor_name = $scope.aq_state.sensors[id].name;
    $scope.aq_form = 'sensor';
  };

  $scope.editRelay = function (id) {
    $scope.aq_edit_relay_id = id;
    $scope.aq_edit_relay_name = $scope.aq_state.relays[id].name;
    $scope.aq_form = 'relay';
  };

  $scope.editTask = function (id) {
    $scope.aq_edit_task_id = id;
    if (id >= 0) {
      let relay = $scope.aq_state.relays[$scope.aq_edit_relay_id];
      $scope.aq_edit_task = angular.copy(relay.tasks[id]);
    } else {
      $scope.aq_edit_task = { start: 0, stop: 0, weekdays: 127 };
    }
    $scope.aq_form = 'task';
  };

  $scope.weekDays = function () {
    return weekDays;
  };

  $scope.isWeekDaySelected = function (day) {
    return $scope.aq_edit_task.weekdays & weekDayBit(day);
  };

  $scope.toggleWeekDay = function (day) {
    $scope.aq_edit_task.weekdays ^= weekDayBit(day);
  };

  $scope.showState();
  $interval($scope.getState, 30000);
});