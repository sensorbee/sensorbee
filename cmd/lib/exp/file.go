package exp

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func setUpFile() cli.Command {
	cmd := cli.Command{
		Name:   "file",
		Usage:  "get a filename of a cache file of a specific data source node",
		Action: runFile,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "all, a",
				Usage: "show all cache files of a node",
			},
			// TODO: --hash: an argument is a hash value rather than a node name
			// TODO: --states: show filenames of states at the time of the cache
		},
	}
	return cmd
}

func runFile(c *cli.Context) error {
	err := func() error {
		if len(c.Args()) != 1 {
			cli.ShowSubcommandHelp(c)
			os.Exit(1)
		}

		cache, err := LoadCacheFromFile(cacheInfoFilename)
		if err != nil {
			return fmt.Errorf("cannot load cache information:%v", err)
		}

		node := c.Args()[0]
		ent := cache.Node(node)
		if ent != nil {
			fmt.Println(CacheFilename(ent))
		}

		if !c.Bool("all") {
			if ent == nil {
				return fmt.Errorf("no cache entry was found for '%v'", node)
			}
			return nil
		}

		found := false
		for _, e := range cache.Outdated {
			if e.NodeName == node {
				found = true
				fmt.Println(CacheFilename(e)) // TODO: sort entries by their timestamps
			}
		}
		if !found {
			return fmt.Errorf("no cache entry was found for '%v'", node)
		}
		return nil
	}()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}
