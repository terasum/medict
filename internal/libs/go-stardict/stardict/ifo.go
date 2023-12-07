package stardict

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// RawInfo contains dictionary options
type RawInfo struct {
	Version string
	Is64    bool
	Options map[string]string
}

func decodeOption(str string) (key string, value string, err error) {
	a := strings.Split(str, "=")

	if len(a) < 2 {
		return "", "", errors.New("Invalid file format: " + str)
	}

	return a[0], a[1], nil
}

// readInfo reads ifo file and collects dictionary options
func readInfo(filename string) (info *RawInfo, err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return
	}

	defer reader.Close()

	r := bufio.NewReader(reader)

	_, err = r.ReadString('\n')

	if err != nil {
		return
	}

	version, err := r.ReadString('\n')

	if err != nil {
		return
	}

	kn, kv, err := decodeOption(version[:len(version)-1])

	if err != nil {
		return
	}

	if kn != "version" {
		err = errors.New("version missing (should be on second line)")
		return
	}

	if kv != "2.4.2" && kv != "3.0.0" && kv != "2.4.2\r" && kv != "3.0.0\r" { // \r : fix the problem on windows
		err = errors.New("stardict version should be either 2.4.2 or 3.0.0")
		return
	}

	info = new(RawInfo)

	info.Version = kv

	info.Options = make(map[string]string)

	for {
		option, err := r.ReadString('\n')

		if err != nil && err != io.EOF {
			return info, err
		}

		if err == io.EOF && len(option) == 0 {
			break
		}

		kn, kv, err = decodeOption(option[:len(option)-1])

		if err != nil {
			return info, err
		}

		info.Options[kn] = kv

		if err == io.EOF {
			break
		}
	}

	if val, ok := info.Options["idxoffsetbits"]; ok {
		if val == "64" {
			info.Is64 = true
		}
	} else {
		info.Is64 = false
	}

	return
}
