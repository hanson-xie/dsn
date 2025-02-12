package app

import (
	"bytes"
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app/config"
	"github.com/Bedrock-Technology/Dsn/log"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
)

type ConfigStore struct {
	mu        sync.RWMutex
	config    map[string]interface{}
	YamlDir   string
	GitUpdate string
	Reload    bool
	Cron      *cron.Cron
}

type ParamType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var (
	instance *ConfigStore
	once     sync.Once
)

func GetConfigStore() *ConfigStore {
	once.Do(func() {
		instance = &ConfigStore{
			config: make(map[string]interface{}),
		}
		log.Info("ConfigStore initialized.")
	})
	return instance
}

func (cs *ConfigStore) LoadConfigs(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".toml" {
			if err := cs.loadFile(path); err != nil {
				log.Infof("Failed to load config %s: %v", path, err)
			}
		}
		return nil
	})
	return err
}

func (cs *ConfigStore) loadFile(path string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var fileConfig map[string]interface{}
	if err := toml.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse TOML file %s: %w", path, err)
	}

	for key, value := range fileConfig {
		cs.config[key] = value
	}
	log.Debugf("load toml", "loadFile: ", path, "content", fileConfig)
	return nil
}

func (cs *ConfigStore) LoadFile(path string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var fileConfig map[string]interface{}
	if err := toml.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse TOML file %s: %w", path, err)
	}

	for key, value := range fileConfig {
		cs.config[key] = value
	}
	log.Infof("load toml", "loadFile: ", path, "content", fileConfig)
	return nil
}

func (cs *ConfigStore) WatchConfigs(dir string, triggerCh <-chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 && filepath.Ext(event.Name) == ".toml" {
					log.Infof("fsnotify monitor", "Detected change in config: ", event.Name)
					_ = cs.loadFile(event.Name)
				}
				if event.Op&fsnotify.Remove != 0 {
					cs.mu.Lock()
					delete(cs.config, event.Name)
					cs.mu.Unlock()
					log.Infof("fsnotify monitor", "Removed config: ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("fsnotify monitor", "Watcher error: ", err)
			case <-triggerCh:
				log.Warn("exit fsnotify monitor")
				return
			}
		}
	}()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Info("Watching for changes in toml files...")
	select {
	case <-triggerCh:
		log.Warn("received sql config shutdown")
		return nil
	}
}

func (cs *ConfigStore) GetDataByKey(key string) (interface{}, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	value, exists := cs.config[key]
	if !exists {
		return nil, fmt.Errorf("func name %s not found", key)
	}
	return value, nil
}

func (cs *ConfigStore) extractVariables(query string) []string {
	// Define a regular expression to match variables prefixed with '@'
	re := regexp.MustCompile(`@([a-zA-Z0-9_\.]+)`)

	// Find all matches
	matches := re.FindAllStringSubmatch(query, -1)

	// Extract the variable names
	var variables []string
	for _, match := range matches {
		if len(match) > 1 {
			variables = append(variables, match[1]) // match[1] contains the variable name
		}
	}
	return variables
}

// BindParams replaces placeholders in the query with corresponding values from paramsValues.
// This function directly replaces placeholders without type validation.
func (cs *ConfigStore) BindParams(query string, paramsJson string) (string, error) {
	paramsNames := cs.extractVariables(query)
	paramsValues := make(map[string]string)
	for _, name := range paramsNames {
		value := gjson.Get(paramsJson, name)
		if !value.Exists() {
			return "", fmt.Errorf("missing value for key %s", name)
		}
		paramsValues[name] = value.String()
	}
	for key, value := range paramsValues {
		placeholder := "@" + key
		// For SQL, wrap string values in single quotes
		//if strings.Contains(value, "'") {
		//	return "", fmt.Errorf("invalid value for key %s: contains single quote", key)
		//}
		query = strings.ReplaceAll(query, placeholder, value)
	}
	return query, nil
}

func (cs *ConfigStore) GetFucList() ([]map[string]string, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	var funcList []map[string]string
	for key := range cs.config {
		content, exists := cs.config[key]
		if !exists {
			return nil, fmt.Errorf("func name %s not found", key)
		}
		contentMap := content.(map[string]interface{})
		params := contentMap[config.SqlParams]
		if params != nil {
			funcList = append(funcList, map[string]string{"func_name": key, "params": params.(string)})
		} else {
			funcList = append(funcList, map[string]string{"func_name": key})
		}
	}
	return funcList, nil
}

func (cs *ConfigStore) IsSafeInput(input string) bool {
	unsafeKeywords := map[string]struct{}{
		";":      {},
		"--":     {},
		"/*":     {},
		"*/":     {},
		"OR":     {},
		"AND":    {},
		"=":      {},
		"DROP":   {},
		"INSERT": {},
		"UPDATE": {},
		"DELETE": {},
	}
	for keyword := range unsafeKeywords {
		if strings.Contains(input, keyword) {
			return false
		}
	}
	return true
}

func (cs *ConfigStore) RunWithWatch(configDir string, triggerCh <-chan struct{}) {
	if err := cs.LoadConfigs(configDir); err != nil {
		panic("Failed to load configs: " + err.Error())
	}

	go func() {
		if err := cs.WatchConfigs(configDir, triggerCh); err != nil {
			log.Errorf("toml watcher", "Failed to watch configs: ", err)
		}
	}()
}

func (cs *ConfigStore) RunWithCron(configDir, gitUpdate, spec string, reload bool, triggerCh <-chan struct{}) {
	if err := cs.LoadConfigs(configDir); err != nil {
		panic("Failed to load configs: " + err.Error())
	}

	cs.YamlDir = configDir
	cs.GitUpdate = gitUpdate
	cs.Reload = reload
	cs.Cron = cron.New(cron.WithLocation(time.UTC), cron.WithSeconds())
	wrappedJob := cron.NewChain(cron.SkipIfStillRunning(cron.DiscardLogger)).Then(cs)
	_, err := cs.Cron.AddJob(spec, wrappedJob)
	if err != nil {
		panic("add job error" + err.Error())
	}

	go cs.start(triggerCh)
}

func (cs *ConfigStore) Run() {
	if cs.Reload {
		log.Info("reload flag is true, start to reload sql config")
		cs.reloadSqlConfig()
	}
}

func (cs *ConfigStore) reloadSqlConfig() {
	cmd := exec.Command("/bin/bash", cs.GitUpdate)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.Errorf("update git", "Failed to execute script:", err)
		log.Errorf("update git", "Standard Error:", stderr.String())
		return
	} else {
		log.Infof("update git,Execution Successful")
		log.Infof("update git,Standard Output:", stdout.String())
	}
	if err := cs.LoadConfigs(cs.YamlDir); err != nil {
		log.Errorf("Failed to load configs: %v", err)
	}
}

func (cs *ConfigStore) start(stopCh <-chan struct{}) {
	cs.Cron.Start()
	log.Info("dns sql center is working")
	select {
	case <-stopCh:
		log.Warn("BedRock dns sql received shutdown")
		contextStop := cs.Cron.Stop()
		select {
		case <-contextStop.Done():
			return
		}
	}
}
