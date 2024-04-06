package main

import (
	"ginchat/models"
	"ginchat/router"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	InitTimer()
	r := gin.Default()
	router.Router(r)
	r.Run(":8080")
}

func InitTimer() {
	utils.Timer(
		time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second,
		time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second,
		models.CleanConnection, "")
}
