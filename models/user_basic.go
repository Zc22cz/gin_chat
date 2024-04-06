package models

import (
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	Salt          string
	LoginTime     time.Time
	HeartBeatTime time.Time
	LoginOutTime  time.Time
	IsLoginOut    bool
	DeviceInfo    string
}
