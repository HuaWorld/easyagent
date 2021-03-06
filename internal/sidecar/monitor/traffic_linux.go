package monitor

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"easyagent/internal/sidecar/base"
)

var (
	re = regexp.MustCompile(fmt.Sprintf(`(?s)class htb %s:(\d+).+?Sent (\d+) bytes`, base.Major))

	cmMu     sync.Mutex
	classMap = map[uint32]uint64{}
)

func getTraffic(classid uint32) (uint64, uint64, error) {
	if classid == 0 {
		return 0, 0, nil
	}

	cmMu.Lock()
	defer cmMu.Unlock()
	return classMap[classid], 0, nil
}

func setTrafficEnable(pid uint32) error { return nil }

func tcStatistic() {
	major, _ := strconv.ParseUint(base.Major, 16, 16)
	major <<= 16

	for {
		time.Sleep(monitorInterval)

		classMapTmp := make(map[uint32]uint64)
		for _, dev := range base.GetDevs() {
			out, err := exec.Command("tc", "-s", "class", "show", "dev", dev).CombinedOutput()
			if err != nil {
				if len(out) > 0 {
					base.Errorf("tc show %v statistic fail: %q", dev, out)
				} else {
					base.Errorf("tc show %v statistic fail: %v", dev, err)
				}
				continue
			}

			for _, match := range re.FindAllSubmatch(out, -1) {
				minor, _ := strconv.ParseUint(string(match[1]), 16, 16)
				sent, _ := strconv.ParseUint(string(match[2]), 10, 0)
				classMapTmp[uint32(major+minor)] += sent
			}
		}
		cmMu.Lock()
		classMap = classMapTmp
		cmMu.Unlock()
	}
}
