package main

import (
	"log"
	"net"
	"net/http"
	"runtime"

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
			log.Printf("%s %s %s %s %s\n", dev.Product(), dev.ProductVersion(), dev.Name(), dev.Location(), dev.UUID())

			// wifictrl request succeeded HTTP 200 OK
			resp, err := http.Get("http://<sonos_ip>:1400/wifictrl?wifi=persist-off")
			defer resp.Body.Close()
			if err != nil {
				log.Println(err)
			}
			if resp.StatusCode == 200 {
				log.Println("Disabled WiFi for", dev.Name())
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
