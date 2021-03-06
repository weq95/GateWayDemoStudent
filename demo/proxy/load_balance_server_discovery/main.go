package main

import (
	"GateWayDemoStudent/demo/proxy/middleware"
	"GateWayDemoStudent/proxy/load_balance"
	proxy2 "GateWayDemoStudent/proxy/proxy"
	"log"
	"net/http"
)

var addr = "127.0.0.1:2002"

func main() {
	mConf, err := load_balance.NewLoadBalanceZkConf("http://%s/base",
		"real_server",
		[]string{"127.0.0.1:2181"},
		map[string]string{"127.0.0.1:2003": "20"})

	if err != nil {
		panic(err)
	}

	rb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbWeightRoundRobin, mConf)

	proxy := proxy2.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)

	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
