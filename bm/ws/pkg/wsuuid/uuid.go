// go get github.com/sony/sonyflake/v2
package wsuuid

import (
	"errors"
	"time"

	"github.com/sony/sonyflake"
)

var uuidInfo *UuidInfo

type UuidInfo struct {
	sonyflake *sonyflake.Sonyflake
	machineId uint16
}

func NewUuidInstance(machineId uint16) error {
	setting := sonyflake.Settings{
		StartTime: time.Now().Add(-time.Hour * 24),
		MachineID: func() (uint16, error) {
			return machineId, nil
		},
	}

	if uuidInfo == nil {
		ins := sonyflake.NewSonyflake(setting)
		if ins == nil {
			return errors.New("创建雪花算法实例失败")
		}

		uuidInfo = &UuidInfo{
			sonyflake: ins,
			machineId: machineId,
		}

	}

	return nil
}

func GetUUID() (uint64, error) {
	uuid, err := uuidInfo.sonyflake.NextID()
	if err != nil {
		return 0, err
	}

	return uuid, nil
}
