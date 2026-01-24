package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"github.com/skypbc/goutils/gerrors"
	"github.com/skypbc/goutils/gnum"
	"strings"
	"sync"
)

func IsLeafDict(m map[string]any) bool {
	for _, v := range m {
		_, ok := v.([]any)
		if !ok {
			return false
		}
	}
	return true
}

func IsFlatIntMap(m map[string]any) bool {
	for _, v := range m {
		_, ok := gnum.TryInt(v)
		if !ok {
			return false
		}
	}
	return true
}

func SanitizeName(name string) string {
	replacements := map[string]string{
		":": "_",
		"+": "_",
	}
	sanitized := name
	for from, to := range replacements {
		sanitized = strings.ReplaceAll(sanitized, from, to)
	}
	return strings.ToUpper(sanitized)
}

func JoinPascal(parts ...string) string {
	var out string
	for _, part := range parts {
		for _, p := range strings.Split(part, "_") {
			p = strings.ToLower(p)
			out += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return out
}

func GetMap(m map[string]any, key string) map[string]any {
	parts := strings.Split(key, ".")
	current := m
	for _, part := range parts {
		if val, ok := current[part]; ok {
			if m, ok := val.(map[string]any); ok {
				current = m
			} else if part == parts[len(parts)-1] {
				return nil
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return current
}

func Map2Struct(m any, s any) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, s)
	if err != nil {
		return err
	}
	return nil
}

func ZipFiles(files map[string]string) (data []byte, err error) {
	var buff bytes.Buffer
	zipWriter := zip.NewWriter(&buff)

	closeOnce := func() func() error {
		var once sync.Once
		var closeErr error
		return func() error {
			once.Do(func() {
				closeErr = zipWriter.Close()
			})
			return closeErr
		}
	}()
	defer func() {
		if err2 := closeOnce(); err2 != nil {
			err = gerrors.Wrap(err, err2).
				SetTemplate("failed to close zip writer")
		}
	}()

	for filePath, content := range files {
		data := []byte(content)
		zipEntry, err := zipWriter.Create(filePath)
		if err != nil {
			return nil, gerrors.Wrap(err).
				SetTemplate("failed to create zip entry")
		}
		if _, err := zipEntry.Write(data); err != nil {
			return nil, gerrors.Wrap(err).
				SetTemplate("failed to write to zip entry")
		}
	}

	if err = closeOnce(); err != nil {
		return nil, nil
	}

	if data = buff.Bytes(); err != nil {
		return nil, gerrors.Wrap(err).
			SetTemplate("failed to get zip buffer")
	}

	return data, nil
}
