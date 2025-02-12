package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"time"

	"github.com/expzhizhuo/fingerprintx/pkg/plugins"
	"github.com/expzhizhuo/fingerprintx/pkg/scan"
)

func main() {
	// setup the scan config (mirrors command line options)
	fxConfig := scan.Config{
		DefaultTimeout: time.Duration(2) * time.Second,
		FastMode:       false,
		Verbose:        false,
		UDP:            false,
	}

	// create a target list to scan
	ip, _ := netip.ParseAddr("220.172.40.182")
	target := plugins.Target{
		Address: netip.AddrPortFrom(ip, 5236),
		Host:    "220.172.40.182",
	}
	targets := make([]plugins.Target, 1)
	targets = append(targets, target)

	// run the scan
	results, err := scan.ScanTargets(targets, fxConfig)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	// process the results
	for _, result := range results {
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatalf("JSON 编码失败: %s\n", err)
		}
		fmt.Println(string(jsonData))
		fmt.Printf("%s:%d (%s/%s)\n", result.Host, result.Port, result.Transport, result.Protocol)
	}
}
