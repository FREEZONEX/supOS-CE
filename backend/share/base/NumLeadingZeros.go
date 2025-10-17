package base

func NumberOfLeadingZerosL(i int64) int {
	var x = int(uint64(i) >> 32)
	if x == 0 {
		return 32 + NumberOfLeadingZeros(int(i))
	}
	return NumberOfLeadingZeros(x)
}

func NumberOfLeadingZeros(i int) int {
	// HD, Count leading 0's
	if i == 0 {
		return 32
	} else if i < 0 {
		return 0
	}

	var n = 31
	if i >= 1<<16 {
		n -= 16
		i = int(uint32(i) >> 16)
	}
	if i >= 1<<8 {
		n -= 8
		i = int(uint32(i) >> 8)
	}
	if i >= 1<<4 {
		n -= 4
		i = int(uint32(i) >> 4)
	}
	if i >= 1<<2 {
		n -= 2
		i = int(uint32(i) >> 2)
	}
	return n - int(uint32(i)>>1)
}
