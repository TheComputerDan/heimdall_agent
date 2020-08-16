package host

import (
	_ "fmt"
	"github.com/spf13/viper"
	"net"
	"os"
	"runtime"
)

// Info defines basic infromation on the Host being collected
// by the agent to uniquely identify the machine in the inventory.
type Info struct {
	IP       map[string]string
	OSType   string
	Hostname string
}

//loadConfig takes agent.yml and loads its values for runtime.
func loadConfig() string {
	viper.SetConfigName("agent.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("docker_agent/config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	interfaceName := viper.GetString("interface_name")
	return interfaceName
}

//BuildInfo returns an `Info` type with the information about the host.
func BuildInfo() Info {
	interfaceName := loadConfig()

	hostInfo := Info{
		IP:       IPAddr(interfaceName),
		OSType:   operatingSystem(),
		Hostname: hostname(),
	}
	return hostInfo
}

// OperatingSystem simply returns `runtime.GOOS`
func operatingSystem() string {
	return runtime.GOOS
}

// Hostname gets the hostname of the device and returns the string
func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

//IPAddr Prints off the IP addresses of each interface
func IPAddr(iface string) map[string]string {
	var addrs map[string]string

	netIface, err := net.InterfaceByName(iface)
	// ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	// fmt.Println(netIface)

	addresses, err := netIface.Addrs()

	osType := operatingSystem()
	if osType == "linux" {
		addrs = map[string]string{"ipv6": addresses[1].String(), "ipv4": addresses[0].String()}
	} else if osType == "darwin" {
		addrs = map[string]string{"ipv6": addresses[0].String(), "ipv4": addresses[1].String()}
	} else {
		addrs = map[string]string{"ipv6": addresses[0].String(), "ipv4": addresses[1].String()}
	}

	// for k, v := range addresses {
	// 	fmt.Printf("Interface Address #%v : %v\n", k, v.String())
	// }
	// return addresses[1].String()
	return addrs
}
