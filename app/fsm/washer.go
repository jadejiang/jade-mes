package fsm

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/transition"

	"jade-mes/app/infrastructure/constant"
	"jade-mes/app/model"
)

// NewWasherStateMachine ...
func NewWasherStateMachine() *StateMachine {
	sm := StateMachine{}
	sm.runner = transition.New(&model.Device{})

	initReserve(sm.runner)
	return &sm
}

// --------------------------------------状态机逻辑---------------------------------------

func initReserve(stateMachine *transition.StateMachine) {
	// 状态：占用中
	stateMachine.State(constant.DeviceStatusRunning).Enter(func(data interface{}, tx *gorm.DB) error {
		return nil
	})

	// 事件：占用（预约）
	stateMachine.Event(constant.DeviceEventOccupy).
		To(constant.DeviceStatusRunning).
		From(constant.DeviceStatusInitial).
		Before(func(data interface{}, tx *gorm.DB) error {
			return nil
		})
}
