package main

import (
	"fmt"
	"net"

	"jiho-dev.com/range-split/iputil"
)

func test1() {
	/*
	    012345678901234567890123456789
	   1   3.....9
	   2         9...............5
	   3   3.....9
	   4   3
	   5            2........1
	   6       7...........9
	   70............................9
	   8     5
	   9      6

	   0....2: 7
	   3....3: 1,3,4,7
	   4....4: 1,3,7
	   5....5: 1,3,7,8
	   6....6: 1,3,7,9
	   7....8: 1,3,6,7
	   9....9: 1,2,3,6,7
	   10..11: 2,6,7
	   12..19: 2,5,6,7
	   20..21: 2,5,7
	   22..25: 2,7
	   26..29: 7
	*/

	var rs RangeSplit

	rs.Init()

	rs.AddRange(3, 9, 1, []uint64{1})
	rs.AddRange(9, 25, 2, []uint64{2})
	rs.AddRange(3, 9, 3, []uint64{3})
	rs.AddRange(3, 3, 4, []uint64{4})
	rs.AddRange(12, 21, 5, []uint64{5})
	rs.AddRange(7, 19, 6, []uint64{6})
	rs.AddRange(0, 29, 7, []uint64{7})
	rs.AddRange(5, 5, 8, []uint64{8})
	rs.AddRange(6, 6, 9, []uint64{9})

	rs.Build()
	rs.DumpIntervals("Done")
}

func test2() {
	/*
	    012345678901234567890123456789
	   1          0
	   2                6
	   3             3.....9

	   10..10: 1
	   13..15: 3
	   16..16: 2,3
	   17..19: 3
	*/

	var rs RangeSplit

	rs.Init()

	rs.AddRange(10, 10, 1, []uint64{1})
	rs.AddRange(16, 16, 2, []uint64{2})
	rs.AddRange(13, 19, 3, []uint64{3})

	rs.Build()
	rs.DumpIntervals("Done")
}

func test3() {
	/*
		80....80: 1
		8080..8080: 2
		8090..8090: 3
	*/

	var rs RangeSplit

	rs.Init()

	rs.AddRange(80, 80, 1, []uint64{1})
	rs.AddRange(8080, 8080, 2, []uint64{2})
	rs.AddRange(8090, 8090, 3, []uint64{3})

	rs.Build()
	rs.DumpIntervals("Done")
}

func test4() {
	/*
		80....80: 1
		8000..8079: 3
		8080..8080: 2, 3
		8081..9000: 3
	*/

	var rs RangeSplit

	rs.Init()

	rs.AddRange(80, 80, 1, []uint64{1})
	rs.AddRange(8080, 8080, 2, []uint64{2})
	rs.AddRange(8000, 9000, 3, []uint64{3})

	rs.Build()
	rs.DumpIntervals("Done")
}

func test5() {
	var rs RangeSplit

	rs.Init()

	ip := net.ParseIP("1.1.1.1")
	s := int64(iputil.Ip2int(ip))
	rs.AddRange(s, s, 1, []uint64{1})
	fmt.Printf("IP: %s, %d \n", ip.String(), s)

	ip = net.ParseIP("1.1.1.10")
	s = int64(iputil.Ip2int(ip))
	rs.AddRange(s, s, 2, []uint64{2})
	fmt.Printf("IP: %s, %d \n", ip.String(), s)

	ip = net.ParseIP("2.1.1.10")
	s = int64(iputil.Ip2int(ip))
	rs.AddRange(s, s, 3, []uint64{3})
	fmt.Printf("IP: %s, %d \n", ip.String(), s)

	rs.Build()
	rs.DumpIntervals("Done")
}

func main() {
	test1()
	test2()
	test3()
	test4()
	test5()

}
