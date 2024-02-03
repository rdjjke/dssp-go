package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	unix, dssp := curTime()
	fmt.Printf("%064s\n", strconv.FormatInt(unix, 2))
	fmt.Printf("%064s\n", strconv.FormatUint(uint64(dssp), 2))

	for {
		_, dssp = curTime()
		fmt.Printf("%v\n", time.Microsecond*time.Duration(dssp))
		time.Sleep(time.Second)
	}
}

func curTime() (unix int64, dssp uint32) {
	unixMicro := time.Now().UnixMicro()
	dsspTime := uint32(unixMicro)
	return unixMicro, dsspTime
}
