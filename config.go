package main

import "sync/atomic"

type Config struct {
	fileserverHits atomic.Int32
}

func (c *Config) GetHits() int {
	return int(c.fileserverHits.Load())
}
