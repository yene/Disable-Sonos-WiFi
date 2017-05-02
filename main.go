package main

import (
	"log"
	"net"
	"net/http"
	"runtime"

	"strings"

	"github.com/ianr0bkny/go-sonos/ssdp"
)

func main() {
	log.Print("Searching for Sonos\n")

	mgr := ssdp.MakeManager()
	mgr.Discover(guessMainInterfaceName(), "11209", false)
	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey("schemas-upnp-org-ContentDirectory"): -1,
	}
	result := mgr.QueryServices(qry)
	if dev_list, has := result["schemas-upnp-org-ContentDirectory"]; has {
		for _, dev := range dev_list {
			location := string(dev.Location())
			address := strings.Replace(location, ":1400/xml/device_description.xml", "", 1)
			ip := strings.Replace(address, "http://", "", 1)
			log.Printf("Found %s %s\n", dev.Product(), ip)

			// wifictrl request succeeded HTTP 200 OK
			resp, err := http.Get(address + ":1400/wifictrl?wifi=persist-off")
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				log.Println("Disabled WiFi for", dev.Product(), ip)
			}
			break
		}
	}
	mgr.Close()
}

func guessMainInterfaceName() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					log.Println("Searching on network adapter", i.Name)
					return i.Name
				}
			}
		}
	}

	if runtime.GOOS == "darwin" {
		return "en1"
	}

	return "eth0"
}
