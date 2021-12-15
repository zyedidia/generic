package hashmap

type Uint64 uint64

func (u Uint64) Hash() uint64 {
	u ^= u >> 33
	u *= 0xff51afd7ed558ccd
	u ^= u >> 33
	u *= 0xc4ceb9fe1a85ec53
	u ^= u >> 33
	return uint64(u)
}
