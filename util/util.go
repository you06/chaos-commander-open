package util

import (
	"context"
	"fmt"

	"github.com/juju/errors"

	//"github.com/ngaut/log"
	"go/build"
	"time"
)

// Version information.
var (
	BuildTS   = "None"
	BuildHash = "None"

	alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	GoGithubPath = build.Default.GOPATH
	pingcap      = fmt.Sprintf("%s/src/github.com/pingcap", GoGithubPath)
	ImagePrefix  = "hub.pingcap.net/schrodinger"
)

// ValidateListenAddr validates listen addr
func ValidateListenAddr(addr string) error {
	return nil
	//if strings.Contains(addr, "0.0.0.0") || strings.HasPrefix(addr, ":") {
	//	return errors.Errorf("invalid addr %s, must specified a public network card", addr)
	//}
	//return nil
}

// PrintInfo prints the octopus version information
func PrintInfo() {
	fmt.Println("Git Commit Hash:", BuildHash)
	fmt.Println("UTC Build Time: ", BuildTS)
}

// Now gets current datetime string
func Now() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

// RetryOnError help work with some unstable components such as ansible
func RetryOnError(ctx context.Context, retryCount int, fn func() error) error {
	var err error
	for i := 0; i < retryCount; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		err = fn()
		if err == nil {
			break
		}

		//log.Error(errors.ErrorStack(err))
		Sleep(ctx, 2*time.Second)
	}

	return errors.Trace(err)
}

// Sleep defines special `sleep` with context
func Sleep(ctx context.Context, sleepTime time.Duration) {
	ticker := time.NewTicker(sleepTime)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		return
	}
}
