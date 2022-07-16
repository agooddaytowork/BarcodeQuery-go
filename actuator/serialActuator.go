package actuator

import (
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"strings"
	"time"
)

type SerialActuator struct {
	PortName       string `json:"portName"`
	PortBaudRate   uint
	port           io.ReadWriteCloser
	actuator1State ActuatorState
	actuator2State ActuatorState
}

var controlString = "${actuator1State}{actuator2State}00#"

func (a *SerialActuator) fromStateToString(state ActuatorState) string {
	if state == OnState {
		return "1"
	}
	return "0"
}

func (a *SerialActuator) sendActuatorStates(actuator1 string, actuator2 string) {
	currentControlString := strings.ReplaceAll(controlString, "{actuator1State}", actuator1)
	currentControlString = strings.ReplaceAll(currentControlString, "{actuator2State}", actuator2)

	_, err := a.port.Write([]byte(currentControlString))
	if err != nil {
		log.Println("Write to port error", err)
		a.port.Close()
		a.InitPort()
	}
}

func (a *SerialActuator) InitPort() {
	options := serial.OpenOptions{
		PortName:        a.PortName,
		BaudRate:        a.PortBaudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	tryingToOpenPort := true
	for tryingToOpenPort {
		port, err := serial.Open(options)
		if err != nil {
			log.Println("port open err", port)
			log.Println("Sleep 5 seconds, then try to open port again")
			time.Sleep(5 * time.Second)
		} else {
			tryingToOpenPort = false
			a.port = port
		}
	}
}

func (a *SerialActuator) setActuator1State(state ActuatorState) {
	if a.actuator1State != state {
		a.actuator1State = state
		a.sendActuatorStates(a.fromStateToString(state), "0")
	}
}

func (a *SerialActuator) SetDuplicateActuatorState(state ActuatorState) {
	log.Println("SetDuplicateActuatorState", state)
	a.setActuator1State(state)
}

func (a *SerialActuator) SetErrorActuatorState(state ActuatorState) {
	log.Println("SetErrorActuatorState", state)
	a.setActuator1State(state)
}

func (a *SerialActuator) SetCameraErrorActuatorState(state ActuatorState) {
	log.Println("SetCameraErrorActuatorState", state)
	a.setActuator1State(state)
}
