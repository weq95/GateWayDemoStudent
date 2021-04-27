package load_balance

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type RandomBalance struct {
	curIndex int
	rss      []string
	//观察主体
	conf LoadBalanceConf
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}

	addr := params[0]
	r.rss = append(r.rss, addr)

	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	r.curIndex = rand.Intn(len(r.rss))

	return r.rss[r.curIndex]
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *RandomBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceZkConf); ok {
		confList := conf.GetConf()
		fmt.Println("update get conf:", confList)

		r.rss = []string{}
		for _, ip := range confList {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		confList := conf.GetConf()
		fmt.Println("Update get conf:", confList)
		r.rss = nil

		for _, ip := range confList {
			_ = r.Add(strings.Split(ip, ".")...)
		}
	}
}
