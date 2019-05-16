// tools project tools.go
package misc

import (
	"os/exec"
	"time"
)

func SleepMilliSecond(millsecond int) {
	time.Sleep(time.Duration(millsecond) * time.Millisecond)
}

func Exec(command string, wait bool, args ...string) error {
	cmd := exec.Command(command, args...)
	err := cmd.Start()
	if wait && err == nil {
		err = cmd.Wait()
	}
	return err
}

func DateTime(pattern string) string {
	return time.Now().Format(pattern)
}

const partten = "2006-01-02 15:04:05"

func Timestamp() string {
	return time.Now().Format(partten)
}

func TimestampByPattern(p string) string {
	return time.Now().Format(p)
}
