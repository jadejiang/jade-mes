package fsm

import (
	"strings"
	"jade-mes/app/model"
	"jade-mes/ecode"

	"github.com/qor/transition"
)

// StateMachine ...
type StateMachine struct {
	runner *transition.StateMachine
}

// Trigger trigger an event
func (fsm StateMachine) Trigger(event string, device *model.Device, persist bool, notes ...string) error {
	nDevice := *device

	nDevice.SetState(nDevice.Status)

	err := fsm.runner.Trigger(event, &nDevice, nil, notes...)

	// 状态不正确
	if err != nil {
		if strings.HasPrefix(err.Error(), "failed to perform event") {
			// if err, exists := ecode.ErrInvalidStateForWasher[event]; exists {
			// 	return err
			// }
			return ecode.ErrInvalidDeviceStatus
		}

		return err
	}

	if persist {
		return nDevice.Update()
	}

	return nil
}
