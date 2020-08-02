package balance

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	curIndex int
	rss      []string
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param illegal")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RandomBalance) Next(key string) (string,error) {
	if len(r.rss) == 0 {
		return "",errors.New("not found node")
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex],nil
}
