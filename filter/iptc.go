package filter

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type iptcTag struct {
	index int
	value string
}

type IPTCInputFilter struct{}

func exifCmd() (string, error) {
	path, err := exec.LookPath("exiftool")
	if err != nil {
		return path, &FatalError{err.Error()}
	}
	return path, nil
}

func runExif(info *ImageInfo, arg ...string) (output []byte, err error) {
	cmdPath, err := exifCmd()
	cmd := exec.Command(cmdPath, arg...)
	var writer io.WriteCloser
	writer, err = cmd.StdinPipe()
	if err == nil {
		go func() {
			writer.Write(info.Buf)
			writer.Close()
		}()

		output, err = cmd.Output()
	}
	if err != nil {
		err = &FatalError{fmt.Sprintf("%v: %v", strings.Join(cmd.Args, " "), err)}
	}
	return output, err
}

func (iptc *IPTCInputFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	output, err := runExif(info, "-T", "-Keywords", "-")
	if err == nil {
		tags := make(map[string]iptcTag)
		for i, tag := range strings.Split(string(output), ",") {
			tokens := strings.Split(tag, ":")
			if len(tokens) == 2 {
				info.Properties[tokens[0]] = strings.TrimSpace(tokens[1])
				tags[tokens[0]] = iptcTag{i, strings.TrimSpace(tokens[1])}
			}
		}
	}
	return info, err
}

type IPTCOutputFilter struct{}

func (iptc *IPTCOutputFilter) Process(info *ImageInfo) (*ImageInfo, error) {
	args := make([]string, 0)
	for key, v := range info.Properties {
		args = append(args, fmt.Sprintf("-Keywords=%v:%v", key, v))
	}
	args = append(args, "-o", "-", "-")
	output, err := runExif(info, args...)
	if err == nil {
		info.Buf = output
	}
	return info, err
}
