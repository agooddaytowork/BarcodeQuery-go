package actuator

func GetActuator(actuatorType string, actuatorURI string) BarcodeActuator {

	switch actuatorType {
	case "console":
		return &ConsoleActuator{}

	case "serial":
		return &SerialActuator{
			PortName:       actuatorURI,
			PortBaudRate:   9600,
			actuator1State: Offstate,
			actuator2State: Offstate,
		}
	}

	panic("Not supported actuator type")
}
