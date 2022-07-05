package gtp5gnl

type OID []uint64

func (oid OID) SEID() (uint64, bool) {
	switch len(oid) {
	case 0, 1:
		return 0, false
	default:
		return oid[0], true
	}
}

func (oid OID) ID() (int, bool) {
	switch len(oid) {
	case 0:
		return 0, false
	case 1:
		return int(oid[0]), true
	default:
		return int(oid[1]), true
	}
}

func (oid OID) Equal(a OID) bool {
	n := len(oid)
	m := len(a)
	if n != m {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != oid[i] {
			return false
		}
	}
	return true
}
