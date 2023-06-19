package orderCheck

import (
	"github.com/gin-gonic/gin"
	"github.com/zlilemon/gin_auto/app/billing"
	"github.com/zlilemon/gin_auto/app/device"
	"github.com/zlilemon/gin_auto/pkg/log"
	"time"
)

type IService interface {
	CheckDeviceStatus()
}

type Service struct {
	//repo IRepository
}

var OrderCheckService = new(Service)

func (s *Service) SCheckDeviceStatus() {
	log.Infof("SCheckDeviceStatus - ")

	// 获取支付订单中，当前时间状态下到时间的订单
	var c *gin.Context
	var req billing.BillingStatusCheckReq
	req.CheckUnixTime = time.Now().Unix()
	deviceStatusList, err := billing.BillingService.SGetOrderStatusCheck(c, req)
	if err != nil {
		log.Errorf("SCheckDeviceStatus error, errMsg:%+v", err)
		return
	}

	if len(deviceStatusList) == 0 {
		// 没有需要把设备关闭的记录
		log.Infof("NoDeviceNeedToClosed, current_time:%d", req.CheckUnixTime)
		return
	}

	//有设备需要进行关闭操作
	for _, v := range deviceStatusList {
		// 找到对应的设备, 触发设备的操作（关闭动作）
		operationReq := device.OperationReq{}
		operationResp := device.OperationResp{}

		operationReq.OutTradeNo = v.OutTradeNo
		operationReq.OpenId = v.OpenId
		operationReq.StoreId = v.StoreId
		operationReq.SeatId = v.SeatId
		operationReq.Cmd = "turn_off"

		err := device.DeviceService.SDeviceOperation(c, operationReq, &operationResp)

		if err != nil {
			log.Errorf("failed to turn off device, out_trade_no:%s, openid:%s, "+
				"store_id:%s, seat_id:%s, err:%s", operationReq.OutTradeNo, operationReq.OpenId,
				operationReq.StoreId, operationReq.SeatId,
				operationReq.Cmd, err.Error())
			return
		}
	}

	log.Infof("SCheckDeviceStatus finish, current_unix_time:%d", req.CheckUnixTime)
}
