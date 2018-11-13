package blocklist

import (
	"fmt"
	"sort"
)

type BlockDesc struct {
	Name    string
	Version int
}

type BlockDescList map[BlockDesc]int

func MakeBlockDescList() BlockDescList {
	return make(map[BlockDesc]int)
}

func (bl BlockDescList) Strings(prefix string) []string {
	var res = make([]string, len(bl))
	i := 0
	for bd, _ := range bl {
		res[i] = fmt.Sprintf("%s%s (v%d)", prefix, bd.Name, bd.Version)
		i++
	}
	sort.Strings(res)
	return res
}

func (bl BlockDescList) Add(name string, version int) {
	bd := BlockDesc{Name: name, Version: version}
	bl[bd] += 1
}
