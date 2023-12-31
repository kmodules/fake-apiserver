package ioutil

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ReadFileAs(path string, obj interface{}) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, obj)
	if err != nil {
		return err
	}
	return nil
}

/*
ReadINIConfig loads a ini config file without any sections. Example:
--- --- ---
a=b
c=d
--- --- ---
*/
func ReadINIConfig(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mp := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result := strings.Split(scanner.Text(), "=")
		if len(result) != 2 {
			continue
		}
		mp[string(result[0])] = result[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mp, nil
}

func WriteJson(path string, obj interface{}) error {
	d, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	EnsureDirectory(path)
	err = os.WriteFile(path, d, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func WriteString(path string, data string) bool {
	EnsureDirectory(path)
	err := os.WriteFile(path, []byte(data), os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

func AppendToFile(path string, values string) error {
	EnsureDirectory(path)
	if _, err := os.Stat(path); err != nil {
		os.WriteFile(path, []byte(""), os.ModePerm)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("\n")
	_, err = f.WriteString(values)
	if err != nil {
		return err
	}
	return nil
}

func EnsureDirectory(path string) error {
	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); err != nil {
		return os.MkdirAll(parent, os.ModePerm)
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func PathNotExists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, fs.ErrNotExist)
}

// WriteFile writes the contents from src to dst using io.Copy.
// If dst does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.
func WriteFile(dst string, src io.Reader, perm os.FileMode) (err error) {
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, src)
	return
}
