package configparam

import (
	"os"
	"strings"
)

func ExtractConfig() string {
	for i, arg := range os.Args {
		if arg == "-config" || arg == "-c" {
			if i+1 < len(os.Args) {
				return os.Args[i+1]
			}
		} else if strings.HasPrefix(arg, "-config=") {
			return strings.TrimPrefix(arg, "-config=")
		} else if strings.HasPrefix(arg, "-c=") {
			return strings.TrimPrefix(arg, "-c=")
		}
	}
	if envCfg := os.Getenv("CONFIG"); envCfg != "" {
		return envCfg
	}
	return ""
}
