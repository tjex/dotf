package util

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Expands environment variables and `~`, returning an absolute path.
func ExpandPath(path string) (string, error) {

	if strings.HasPrefix(path, "~") {
		// resolve necessary variables
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		home := usr.HomeDir
		if home == "" {
			err := errors.New("Could not find user's home directory.")
			return "", err
		}

		// construct as abs path
		if path == "~" {
			path = home
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home, path[2:])
		}
	}

	path = os.ExpandEnv(path)

	return path, nil
}



