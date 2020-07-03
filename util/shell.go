package util

import (
	"context"
	"io"
	"os/exec"
)

func DoCmd(dir string, c string, args ...string) (string, error) {
	cmd := exec.Command(c, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func DoCmdContext(ctx context.Context, dir string, c string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, c, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func DoCmdContextWithLogger(ctx context.Context, writer io.Writer, dir string, c string, args ...string) error {
	cmd := exec.CommandContext(ctx, c, args...)
	cmd.Dir = dir
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
