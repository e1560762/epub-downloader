package main

import (
	"e1560762/urldownloader/pkg"
	"flag"
	"fmt"
	"sync"
)

func main() {
	configFile := flag.String("config", "config.json", "Path to th config file")
	sourceURLS := pkg.NewOriginURLsArg()
	flag.Var(&sourceURLS, "remote-urls", "List of source URLs to download from")
	filePaths := pkg.NewFilePathsArg()
	flag.Var(&filePaths, "file-paths", "List of paths to download the files")
	flag.Parse()

	if len(sourceURLS) != len(filePaths) {
		fmt.Printf("[ERROR] number of urls does not match with the number of file paths")
		return
	}
	conf, err := pkg.NewConfig(*configFile)
	if err != nil {
		fmt.Printf("[ERROR] Parsing %s: %v\n", *configFile, err)
		return
	}
	if conf == nil || len(conf.HostConf) == 0 {
		fmt.Printf("[ERROR] Config file is empty %s: %v\n", *configFile, err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(sourceURLS))
	for i := 0; i < len(sourceURLS); i++ {
		go func(c pkg.Config, u, p string) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[PANIC] Downloading %s to %s: %v\n", u, p, r)
					wg.Done()
				}
			}()
			if err = pkg.NewParser(c, u, p).Download(); err != nil {
				fmt.Printf("[ERROR] Downloading %s to %s: %v\n", u, p, err)
			}
			wg.Done()
		}(*conf, sourceURLS[i], filePaths[i])
	}
	wg.Wait()
}
