package dsn

import (
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app/config"
	"github.com/Bedrock-Technology/Dsn/log"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	conf *config.DsnConfig
	mu   sync.RWMutex
)

func LoadConfig(confPath string) error {
	filename, err := filepath.Abs(confPath)
	if err != nil {
		log.Errorf("LoadConfig", "Failed to get absolute path of config file:", err)
		return err
	}
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		log.Errorf("LoadConfig", "Config file does not exist:", filename)
		return err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if conf == nil {
		conf = &config.DsnConfig{}
	}

	if err := yaml.Unmarshal(data, conf); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}

func RunWithWatch(confPath string, triggerCh <-chan struct{}) {
	filename, err := filepath.Abs(confPath)
	if err != nil {
		log.Errorf("RunWithWatch", "Failed to get absolute path of config file:", err)
		return
	}
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		log.Errorf("RunWithWatch", "Config file does not exist:", filename)
		return
	}
	go func() {
		if err := watchConfig(filename, triggerCh); err != nil {
			log.Errorf("RunWithWatch", "Failed to watch configs:", err)
		}
	}()
}

func watchConfig(filename string, triggerCh <-chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	dir := filepath.Dir(filename)
	if err := watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to add directory to watch list: %v", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Error("fsnotify monitor: event channel closed")
					return
				}
				if (event.Name == filename || filepath.Base(event.Name) == filepath.Base(filename)) && event.Op&(fsnotify.Write) != 0 {
					log.Infof("fsnotify monitor", "Detected event: %s on file: %s", event.Op, event.Name)
					time.Sleep(200 * time.Millisecond)
					if _, err := os.Stat(filename); err == nil {
						if err := LoadConfig(filename); err != nil {
							log.Errorf("fsnotify monitor", "Failed to reload config: %v", err)
						} else {
							log.Infof("fsnotify monitor", "Successfully reloaded config file")
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Errorf("fsnotify monitor", "Watch error: ", err)
					return
				}
			case <-triggerCh:
				log.Info("fsnotify monitor Exit fsnotify monitor")
				watcher.Remove(dir)
				return
			}
		}
	}()

	log.Infof("fsnotify monitor", "Watching config file and directory: %s", dir)
	select {
	case <-triggerCh:
		log.Info("received yaml config shutdown")
		return nil
	}
	return nil
}

func GetConfig() *config.DsnConfig {
	mu.RLock()
	defer mu.RUnlock()
	return conf
}
