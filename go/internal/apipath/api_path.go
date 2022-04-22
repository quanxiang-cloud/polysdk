package apipath

import (
	"path"
	"strings"
)

// Split parse full path with path and name
func Split(full string) (string, string) {
	if !strings.HasPrefix(full, "/") {
		full = "/" + full
	}
	path, name := "", full
	if index := strings.LastIndex(full, "/"); index >= 0 {
		path = full[:index]
		name = full[index+1:]
	}
	return path, name
}

// Name get the last name of a full path
func Name(fullPath string) string {
	_, name := Split(fullPath)
	return name
}

// Parent get the parent path of a full path
func Parent(fullPath string) string {
	parent, _ := Split(fullPath)
	return parent
}

// Join join the namespace and name as full path
func Join(namespace, name string) string {
	if !strings.HasPrefix(namespace, "/") {
		namespace = "/" + namespace
	}
	return path.Join(namespace, name)
}

// Format convert full path as standard format
func Format(full string) string {
	path, name := Split(full)
	return Join(path, name)
}
