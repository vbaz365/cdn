package utils

func commonPrefixLen(node uint64, inserted uint64) uint8 {
	var count uint8
	var differentBits uint64 = node ^ inserted
	for differentBits != 0 {
		differentBits = differentBits >> 1
		count++
	}
	return 64 - count
}
