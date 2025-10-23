package bytes

import "sort"

func rawToSet(a []byte, size int) map[string]interface{} {
	r := make(map[string]interface{})
	for i := 0; i < len(a); i += size {
		r[string(a[i:i+size])] = nil
	}
	return r
}

func Union(a []byte, b []byte, size int) []byte {
	aa := rawToSet(a, size)
	bb := rawToSet(b, size)

	r := make(map[string]interface{})
	for k := range aa {
		_, ok := bb[k]
		if ok {
			r[k] = nil
		}
	}
	resp := make([]byte, size*len(r))
	i := 0
	for k := range r {
		copy(resp[i*size:(i+1)*size], []byte(k))
		i++
	}
	return resp
}

func WhatsUp(old []byte, fresh []byte, size int) []byte {
	old_s := rawToSet(old, size)
	todo := make(map[string]interface{})
	var ok bool
	var key string
	for id := 0; id < len(fresh); id += size {
		key = string(fresh[id : id+size])
		_, ok = old_s[key]
		if !ok {
			todo[key] = nil
		}
	}
	r := make(ById, len(todo))
	i := 0
	for k := range todo {
		r[i] = Id(k)
		i++
	}
	sort.Sort(r)
	rr := make([]byte, len(r)*size)
	for i, id := range r {
		for s := range size {
			rr[i*size+s] = id[s]
		}
	}
	return rr
}
