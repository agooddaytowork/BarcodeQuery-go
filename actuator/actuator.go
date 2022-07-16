package actuator

type ActuatorState bool

func GetState(state bool) ActuatorState {
	if state {
		return OnState
	}

	return Offstate
}

const (
	OnState  ActuatorState = true
	Offstate ActuatorState = false
)

type BarcodeActuator interface {
	SetDuplicateActuatorState(state ActuatorState)
	SetErrorActuatorState(state ActuatorState)
	SetCameraErrorActuatorState(state ActuatorState)
}
