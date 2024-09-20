package consul

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"jade-mes/config"
)

// RegisterServer ...
func RegisterServer() {

	settings := config.GetConfig()

	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	content := settings.GetString("server.port")[1:len(settings.GetString("server.port"))]
	portCode, err := strconv.Atoi(content)

	var (
		ip          = localIP()
		port        = portCode
		tags        = config.Config.Consul.Tags
		serviceName = config.Config.Consul.ServiceName
	)

	if ip == "" {
		log.Fatal("consul register ip error")
	}

	var healthCheck = &consulapi.AgentServiceCheck{ // 健康检查
		HTTP:                           fmt.Sprintf("http://%s:%d/health", ip, port),
		Timeout:                        "3s",
		Interval:                       "10s", // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
	}

	r := &consulapi.AgentServiceRegistration{
		ID:      hostName + "-" + strconv.Itoa(portCode), // 服务节点的名称
		Name:    serviceName,                             // 服务名称
		Port:    port,                                    // 服务端口
		Tags:    tags,                                    // tag，可以为空
		Address: localIP(),                               // 服务 IP
		Check:   healthCheck,
	}

	err = client.Agent().ServiceRegister(r)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
