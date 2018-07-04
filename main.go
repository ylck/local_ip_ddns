package main

import (
	"net"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

type Hostinfo struct {
	Hostname string
	Ip       string
}

var ips string

func init() {
	//设置最低loglevel
	log.SetLevel(log.InfoLevel)
}

func getip() *Hostinfo {
	// func getip() []byte {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log.Infof("Hostname: %v", hostname)
	// interfaces, err := net.Interfaces()
	// if err != nil {
	// 	log.Info(err)
	// }
	// nic_list := [...]string{"br0", "en0", "br1"}
	// for _, i := range interfaces {

	// byNameInterface, err := net.InterfaceByName(i.Name)
	// nic_name = br1
	byNameInterface, err := net.InterfaceByName(os.Getenv("nic_name"))
	if err != nil {
		log.Error(err)
	}

	addresses, err := byNameInterface.Addrs()
	log.Info(addresses)

	for _, v := range addresses {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = ipnet.IP.String()
				log.Infof("Ip: %v", ipnet.IP.String())

			}
			break
		}

	}
	// }
	info := &Hostinfo{
		Hostname: hostname,
		Ip:       ips,
	}

	return info
}

func DDNS(new_ip string) {
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}
	// zone_name :  ylck.me
	zoneID, err := api.ZoneIDByName(os.Getenv("zone_name"))
	if err != nil {
		log.Fatal(err)
	}
	// Name: unraid.ylck.me
	foo := cloudflare.DNSRecord{Name: os.Getenv("sld_name") + os.Getenv("zone_name")}
	recs, err := api.DNSRecords(zoneID, foo)
	if err != nil {
		log.Fatal(err)
	}
	var id, ip string

	for _, r := range recs {

		id = r.ID
		ip = r.Content
		log.Printf(id, ip)
	}
	if new_ip != ip {
		// Fetch all records for a zone
		err = api.UpdateDNSRecord(zoneID, id, cloudflare.DNSRecord{Content: new_ip})
		if err != nil {
			log.Fatal(err)
		}
		log.Info("DDNS update")
	}
}

func main() {
	var ip = getip().Ip
	gocron.Every(1).Seconds().Do(DDNS, ip)
	<-gocron.Start()
}
