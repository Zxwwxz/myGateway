package balance

type LbType int

const (
	LbRandom LbType = iota
	LbRound
	LbWeight
	LbHash
)

func LoadBanlanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbHash:
		return NewConsistentHashBanlance(10, nil)
	case LbRound:
		return &RoundBalance{}
	case LbWeight:
		return &WeightBalance{}
	default:
		return &RandomBalance{}
	}
}
