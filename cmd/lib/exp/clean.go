package exp

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func setUpClean() cli.Command {
	cmd := cli.Command{
		Name:   "clean",
		Usage:  "clean up cache files",
		Action: runClean,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "all, a",
				Usage: "clean all cache files",
			},
			// --cache <name> (default .sensorbee-exp-cache)
		},
	}
	return cmd
}

func runClean(c *cli.Context) error {
	err := func() (retErr error) {
		cache, err := LoadCacheFromFile(cacheInfoFilename)
		if err != nil {
			return fmt.Errorf("cannot load cache information:%v", err)
		}
		defer func() {
			if err := SaveCacheToFile(cache, cacheInfoFilename); err != nil {
				if retErr == nil { // don't overwrite the previous error
					retErr = fmt.Errorf("cannot save a cache file: %v", err)
				}
			}
		}()

		if err := removeCacheFiles(cache.Outdated); err != nil {
			return err
		}
		cache.Outdated = nil
		if !c.Bool("all") {
			return nil
		}

		// TODO: ask if the user really wants to clean all cache files.
		if err := removeCacheFiles(cache.Applied); err != nil {
			return err
		}
		cache.Applied = nil
		return nil
	}()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func removeCacheFile(ent *CacheEntry) error {
	if ent.NodeName != "" {
		if err := os.Remove(CacheFilename(ent)); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("cannot remove a cache file '%v': %v", ent.NodeName, err)
			}
		}
	}

	for _, name := range ent.States {
		if err := os.Remove(StateCacheFilename(ent, name)); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("cannot remove a cache file '%v': %v", name, err)
			}
		}
	}
	return nil
}

func removeCacheFiles(ents []*CacheEntry) error {
	for _, ent := range ents {
		if err := removeCacheFile(ent); err != nil {
			return err
		}
	}
	return nil
}
