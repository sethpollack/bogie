package main

import "os"

type Env struct {
}

func (e *Env) Getenv(key string, def ...string) string {
	val := os.Getenv(key)
	if val == "" && len(def) > 0 {
		return def[0]
	}

	return val
}
