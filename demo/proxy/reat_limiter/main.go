package reat_limiter

import (
	"GateWayDemoStudent/demo/proxy/middleware"
	"GateWayDemoStudent/demo/proxy/proxy"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"net/url"
)

//熔断方案
var addr = "127.0.0.1:2002"

func main() {
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003/base"

		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			log.Println(err1)
		}

		rs2 := "http://127.0.0.1:2004/base"
		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			log.Println(err2)
		}

		urls := []*url.URL{url1, url2}

		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}

	log.Println("Starting httpserver at " + addr)

	sliceRouter := middleware.NewSliceRouter()
	sliceRouter.Group("/").Use(middleware.RateLimiter())
	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRouter)

	log.Fatal(http.ListenAndServe(addr, routerHandler))
}

func RateLimiter() func(ctx *middleware.SliceRouterContext) {
	l := rate.NewLimiter(1, 2)

	return func(c *middleware.SliceRouterContext) {
		if !l.Allow() {
			c.Rw.Write([]byte(fmt.Sprintf("rate limit: %v:%v", l.Limit(), l.Burst())))

			c.Abort()
			return
		}

		c.Next()
	}
}
