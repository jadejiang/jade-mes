package metric

import (
	"fmt"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	"github.com/influxdata/influxdb1-client/v2"
	"time"
	"jade-mes/config"
)

var udpClient client.Client

type influxConfig struct {
	Host string
	Port string
}

func initUdpClient() {
	settings := config.GetConfig()

	var cfg influxConfig

	_ = settings.UnmarshalKey("influxdb", &cfg)

	addr := cfg.Host + ":" + cfg.Port

	c, err := client.NewUDPClient(client.UDPConfig{Addr: addr})
	if err != nil {
		fmt.Println("influxdb initUdpClient: ", err.Error())
		return
	}

	udpClient = c
	return
}

func GetClient() client.Client {
	if udpClient == nil {
		initUdpClient()
	}

	return udpClient
}

func Send(measurement string, tags map[string]string) {
	c := udpClient
	// Make client
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: "telegraf", Precision: "us"})
	if err != nil {
		fmt.Println("influxdb NewBatchPoints: ", err.Error())
		return
	}

	fields := map[string]interface{}{
		"count": 1,
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		fmt.Println("influxdb newPoint Error: ", err.Error())
		return
	}

	bp.AddPoint(pt)

	// Write the batch
	err = c.Write(bp)
	if err != nil {
		fmt.Println("influxdb write Error: ", err.Error())
		return
	}

	return
}
