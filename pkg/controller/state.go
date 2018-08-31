package controller

type ControllerState struct {
	Sensors []SensorState `json:"sensors"`
	Relays  []RelayState  `json:"relays"`
	Actions []ActionState `json:"actions"`
}

type SensorState struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Value *float32 `json:"value"`
}

type RelayState struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	On   bool   `json:"on"`
}

type ActionState struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	On   bool   `json:"on"`
}
