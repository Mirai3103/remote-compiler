package isolate

import (
	"fmt"
)

type IsolateCommandBuilder struct {
	agrs []string
}

func NewIsolateCommandBuilder() *IsolateCommandBuilder {
	return &IsolateCommandBuilder{
		agrs: []string{"isolate"},
	}
}

func (icb *IsolateCommandBuilder) WithCGroup() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--cg")
	return icb
}
func (icb *IsolateCommandBuilder) WithBoxID(id int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--box-id="+fmt.Sprint(id))
	return icb
}

func (icb *IsolateCommandBuilder) WithCGroupMemory(mem int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--cg-mem="+fmt.Sprint(mem))
	return icb
}

func (icb *IsolateCommandBuilder) WithCGroupTiming() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--cg-timing")
	return icb
}

func (icb *IsolateCommandBuilder) WithChdir(dir string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--chdir="+dir)
	return icb
}

func (icb *IsolateCommandBuilder) AddDir(dir string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--dir="+dir)
	return icb
}

func (icb *IsolateCommandBuilder) AddEnv(env string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--env="+env)
	return icb
}

func (icb *IsolateCommandBuilder) WithExtraTime(time int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--extra-time="+fmt.Sprint(time))
	return icb
}

func (icb *IsolateCommandBuilder) WithFullEnv() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--full-env")
	return icb
}

func (icb *IsolateCommandBuilder) WithInheritFds() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--inherit-fds")
	return icb
}

func (icb *IsolateCommandBuilder) WithMemory(mem int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--mem="+fmt.Sprint(mem))
	return icb
}

func (icb *IsolateCommandBuilder) WithQuota(blk, ino string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--quota=", blk, ",", ino)
	return icb
}

func (icb *IsolateCommandBuilder) WithShareNet() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--share-net")
	return icb
}

func (icb *IsolateCommandBuilder) WithSilent() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--silent")
	return icb
}

func (icb *IsolateCommandBuilder) WithStackSize(size int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--stack="+fmt.Sprint(size))
	return icb
}

func (icb *IsolateCommandBuilder) WithStderrToStdout() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--stderr-to-stdout")
	return icb
}

func (icb *IsolateCommandBuilder) WithProcesses(max int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--processes="+fmt.Sprint(max))
	return icb
}

func (icb *IsolateCommandBuilder) WithTime(time int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--time="+fmt.Sprint(time))
	return icb
}

func (icb *IsolateCommandBuilder) WithTtyHack() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--tty-hack")
	return icb
}

func (icb *IsolateCommandBuilder) WithVerbose() *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--verbose")
	return icb
}

func (icb *IsolateCommandBuilder) WithWallTime(time int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--wall-time="+fmt.Sprint(time))
	return icb
}

func (icb *IsolateCommandBuilder) WithMaxFileSize(size int) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--fsize="+fmt.Sprint(size))
	return icb
}

func (icb *IsolateCommandBuilder) WithStdinFile(file string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--stdin="+file)
	return icb
}

func (icb *IsolateCommandBuilder) WithStdoutFile(file string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--stdout="+file)
	return icb
}

func (icb *IsolateCommandBuilder) WithStderrFile(file string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--stderr="+file)
	return icb
}

func (icb *IsolateCommandBuilder) WithMetaFile(file string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--meta="+file)
	return icb
}

func (icb *IsolateCommandBuilder) WithRunCommand(command string) *IsolateCommandBuilder {
	icb.agrs = append(icb.agrs, "--run", "--", command)
	return icb
}

func (icb *IsolateCommandBuilder) Build() []string {
	return icb.agrs
}

func (icb *IsolateCommandBuilder) Clone() *IsolateCommandBuilder {
	cloned := make([]string, len(icb.agrs))
	copy(cloned, icb.agrs)
	return &IsolateCommandBuilder{
		agrs: cloned,
	}
}
