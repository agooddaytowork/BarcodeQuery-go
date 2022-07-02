package actuator

type ActuatorState bool

const (
	OnState  ActuatorState = true
	Offstate ActuatorState = false
)

type BarcodeActuator interface {
	SetDuplicateActuatorState(state ActuatorState)
	SetErrorActuatorState(state ActuatorState)
}
