package actuator

import "log"

type ConsoleActuator struct {
}

func (a *ConsoleActuator) SetDuplicateActuatorState(state ActuatorState) {
	log.Println("SetDuplicateActuatorState", state)
}

func (a *ConsoleActuator) SetErrorActuatorState(state ActuatorState) {
	log.Println("SetErrorActuatorState", state)
}

func (a *ConsoleActuator) SetCameraErrorActuatorState(state ActuatorState) {
	log.Println("SetCameraErrorActuatorState", state)
}
