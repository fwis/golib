package ini

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var (
	bComment = []byte{'#'}
	bEmpty   = []byte{}
	bEqual   = []byte{'='}
	bDQuote  = []byte{'"'}
)

func ReadIniConfig(configfile string) (map[string]string, error) {
	cfg := make(map[string]string)
	file, err := os.Open(configfile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)

	for nComment, off := 0, int64(1); ; {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmpty) {
			continue
		}

		off += int64(len(line))

		if bytes.HasPrefix(line, bComment) {
			line = bytes.TrimLeft(line, "#")
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}
		if comment.Len() != 0 {
			comment.Reset()
			nComment++
		}

		val := bytes.SplitN(line, bEqual, 2)
		if bytes.HasPrefix([]byte(strings.TrimSpace(string(val[1]))), bDQuote) {
			val[1] = bytes.Trim([]byte(strings.TrimSpace(string(val[1]))), `"`)
		}

		key := strings.TrimSpace(string(val[0]))
		cfg[key] = strings.TrimSpace(string(val[1]))
	}
	return cfg, nil
}

type IniConf struct {
	mm map[string]string
}

var EIniConfNoKey = errors.New("No ini config key")

func ReadIniConfig1(configfile string) (*IniConf, error) {
	mm, err := ReadIniConfig(configfile)
	if err != nil {
		return nil, err
	}
	return &IniConf{mm: mm}, nil
}

func (m *IniConf) ReadInt(k string) (int, error) {
	if v, ok := m.mm[k]; ok {
		return strconv.Atoi(v)
	} else {
		return 0, EIniConfNoKey
	}
}

func (m *IniConf) ReadInt64(k string) (int64, error) {
	if v, ok := m.mm[k]; ok {
		return strconv.ParseInt(v, 10, 64)
	} else {
		return 0, EIniConfNoKey
	}
}

func (m *IniConf) ReadString(k string) string {
	if v, ok := m.mm[k]; ok {
		return v
	} else {
		return ""
	}
}

func (m *IniConf) ReadFloat64(k string) (float64, error) {
	if v, ok := m.mm[k]; ok {
		return strconv.ParseFloat(v, 64)
	} else {
		return 0, EIniConfNoKey
	}
}
