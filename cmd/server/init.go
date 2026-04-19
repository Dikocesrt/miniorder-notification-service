package main

import (
	"encoding/json"
	"fmt"
	"miniorder-order-service/internal/config"
	"miniorder-order-service/pkg/utils"
	"strings"
)

func RunApplication() {
	svcProperties := config.GetServiceProperties()
	fmt.Println(strings.ToUpper(svcProperties.ServiceName) + " SERVICE")
	asBytesJSON, _ := json.Marshal(svcProperties)
	fmt.Printf("SERVICE CONFIG : '%v'\n", string(asBytesJSON))

	db := utils.InitConnectDatabase(svcProperties.ServiceName, svcProperties.DebugMode, svcProperties.Database)
	redis := utils.InitConnectREDIS(svcProperties.Redis)

	fmt.Println(db, redis)
}
