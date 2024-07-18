package os

import (
	"fmt"
	"os"
)

// DefaultSocketPath returns the default path to the kmd.sock
// it is ${XDG_RUNTIME:/run/user/`id`}/kmd.sock
func DefaultSocketPath() string {
	const sockName = "/kmd.sock"
	def := func() string {
		return fmt.Sprintf("/run/user/%d", os.Getuid()) + sockName
	}
	return getEnvOrFunc("XDG_RUNTIME_DIR", def) + sockName
}

func getEnvOrFunc(key string, def func() string) string {
	if value, defined := os.LookupEnv(key); defined {
		return value
	}
	return def()
}
