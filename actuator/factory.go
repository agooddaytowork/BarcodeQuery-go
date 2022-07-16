package actuator

func GetActuator(actuatorType string, actuatorURI string) BarcodeActuator {

	switch actuatorType {
	case "console":
		return &ConsoleActuator{}

	case "serial":
		actuator := SerialActuator{
			PortName:       actuatorURI,
			PortBaudRate:   9600,
			actuator1State: Offstate,
			actuator2State: Offstate,
		}
		actuator.InitPort()
		return &actuator
	}

	panic("Not supported actuator type")
}
