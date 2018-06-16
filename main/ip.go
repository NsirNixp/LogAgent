package main

import (
	"fmt"
	"net"
)

/**
 * 获取本机器的ip地址
 * @author Charelyz
 */

var ips []string

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("The program can't get the address,[error=%v]", err))
	}

	for _, addr := range addrs {
		// 转换成*net.IPNet，检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	fmt.Printf("Ip=[%v]\n", ips)

}
