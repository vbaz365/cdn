package utils

func CommonPrefixLen(node uint64, inserted uint64) uint8 {
    var count uint8
    var differentBits uint64 = node ^ inserted
    for differentBits != 0 {
        differentBits = differentBits >> 1
        count++
    }
    return 64 - count
}

func GetFirstBit(input uint64) uint8 {
    var ret uint8 = uint8(input >> 63)
    return ret
}

func GetBitsAfterIndex(high uint64, low uint64, index uint8) uint64{
	var ret uint64
	if(index > 63){
		ret = low << (index - 64)
	} else{
		ret = (high << index) | low >> (64 - index)
	}
	return ret
}

func CombineSmallCompressions(first uint64, firstLen uint8, second uint64, secondLen uint8) (uint64, uint8){
	var compression uint64
	compressionLen := firstLen + secondLen
	second = second >> firstLen
	compression = first | second
	return compression, compressionLen
}

func SplitCompression(compression uint64, splitIndex uint8) (uint64, uint64){
	var lowerPart uint64 = compression << splitIndex
	var higherPart = (compression >> (64 - splitIndex)) << (64 - splitIndex)
	return higherPart, lowerPart
}