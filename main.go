package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zlilemon/gin_auto/app/user"
	"github.com/zlilemon/gin_auto/pkg/config"
	"github.com/zlilemon/gin_auto/pkg/database"
	"github.com/zlilemon/gin_auto/pkg/log"
	"github.com/zlilemon/gin_cron/app/orderCheck"
)

func helloCron() {
	fmt.Println("hello gin_cron")
}

func main() {
	fmt.Println("start go cron ...")

	config.InitConf()
	log.Infof("mysql username:%s", config.MysqlOption.Username)

	// 链接数据库实例
	database.ConnectMysql()

	cron := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(nil), cron.Recover(nil)))

	cron.AddFunc("*  *  *  *  *  *", helloCron)
	cron.AddFunc("*  *  *  *  *  *", user.UserService.HelloWorld)
	cron.AddFunc("*  *  *  *  *  *", orderCheck.OrderCheckService.SCheckDeviceStatus)
	//log.Infof("hello")

	cron.Start()

	defer cron.Stop()

	select {}
}
