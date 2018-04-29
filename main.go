package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"strings"

	"github.com/ianr0bkny/go-sonos/ssdp"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	fmt.Println("Searching for Sonos")
	for _, n := range interfaceNames() {
		findAndDisableOn(n)
	}

	fmt.Println("Finished")
	time.Sleep(time.Second * 5)

}

func findAndDisableOn(network string) {
	log.Println("Searching on network adapter", network)
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("Error", r)
		}
	}()

	mgr := ssdp.MakeManager()
	err := mgr.Discover(network, "11209", false)
	if err != nil {
		log.Println(err)
		mgr.Close()
		return
	}

	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey("schemas-upnp-org-ContentDirectory"): -1,
	}
	result := mgr.QueryServices(qry)
	if dev_list, has := result["schemas-upnp-org-ContentDirectory"]; has {
		for _, dev := range dev_list {
			if dev.Product() != "Sonos" {
				continue
			}

			location := string(dev.Location())
			address := strings.Replace(location, ":1400/xml/device_description.xml", "", 1)
			ip := strings.Replace(address, "http://", "", 1)
			url := "http://" + ip + ":1400/wifictrl?wifi=persist-off"
			fmt.Println("Found Sonos, open this URL and disable WiFi by selecting Persist Off", url)
		}
	}
	mgr.Close()
}

func interfaceNames() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var ifs []string

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ifs = append(ifs, i.Name)
				}
			}
		}
	}

	return ifs
}
