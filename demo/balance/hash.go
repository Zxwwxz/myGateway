package balance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)
//hash函数类型
type HashFunc func(data []byte) uint32

//排序
type Uint32Silce []uint32

func (u Uint32Silce) Len()(int) {
	return len(u)
}

func (u Uint32Silce) Less(i,j int)(bool) {
	return u[i] < u[j]
}

func (u Uint32Silce) Swap(i,j int)() {
	u[i],u[j] = u[j],u[i]
}

type HashBanlance struct {
	mux      	sync.RWMutex       //锁
	hashFunc 	HashFunc           //哈希函数
	replicas 	int                //虚拟节点个数
	nodesKeys	Uint32Silce        //虚拟节点hashkey排序列表
	nodesMap 	map[uint32]string  //虚拟节点信息
}

func NewConsistentHashBanlance(replicas int, fn HashFunc) *HashBanlance {
	m := &HashBanlance{
		replicas: replicas,
		hashFunc: fn,
		nodesMap: make(map[uint32]string),
	}
	//默认哈希函数
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// 验证是否为空
func (h *HashBanlance) IsEmpty() bool {
	return len(h.nodesKeys) == 0
}

//用来添加缓存节点
func (h *HashBanlance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param illegal")
	}
	addr := params[0]
	h.mux.Lock()
	defer h.mux.Unlock()
	//增加虚拟节点
	for i := 0; i < h.replicas; i++ {
		//计算hash值
		hashKey := h.hashFunc([]byte(strconv.Itoa(i)+addr))
		h.nodesKeys = append(h.nodesKeys,hashKey)
		h.nodesMap[hashKey] = addr
	}
	//进行排序
	sort.Sort(h.nodesKeys)
	return nil
}

//用来查找下一个节点
func (h *HashBanlance) Next(key string) (string, error) {
	if h.IsEmpty() {
		return "", errors.New("not found node")
	}
	hashKey := h.hashFunc([]byte(key))
	//用二分查找，第一个"服务器hash"值大于"数据hash"值的就是最优"服务器节点"
	index := sort.Search(len(h.nodesKeys), func(i int) bool {
		return h.nodesKeys[i] >= hashKey
	})
	//大于服务器节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if index == len(h.nodesKeys) {
		index = 0
	}
	h.mux.RLock()
	defer h.mux.RUnlock()
	return h.nodesMap[h.nodesKeys[index]], nil
}