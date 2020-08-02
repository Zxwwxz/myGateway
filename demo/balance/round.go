package balance

import "errors"

type RoundBalance struct {
	curIndex int
	rss      []string
}

func (r *RoundBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param illegal")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RoundBalance) Next(key string) (string,error) {
	if len(r.rss) == 0 {
		return "",errors.New("not found node")
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr,nil
}
