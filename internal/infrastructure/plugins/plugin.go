// File: internal/infrastructure/plugins/plugin.go

package plugins

import (
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/amirhossein-jamali/goden-crawler/internal/domain/interfaces"
	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
)

// PluginManager manages the loading and registration of plugins
type PluginManager struct {
	pluginDir       string
	extractors      map[string]interfaces.Extractor
	formatters      map[string]interfaces.FormatterService
	extractorsMutex sync.RWMutex
	formattersMutex sync.RWMutex
	loadedPlugins   map[string]*plugin.Plugin
	pluginsMutex    sync.RWMutex
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginDir string) *PluginManager {
	return &PluginManager{
		pluginDir:     pluginDir,
		extractors:    make(map[string]interfaces.Extractor),
		formatters:    make(map[string]interfaces.FormatterService),
		loadedPlugins: make(map[string]*plugin.Plugin),
	}
}

// LoadPlugins loads all plugins from the plugin directory
func (pm *PluginManager) LoadPlugins() error {
	if pm.pluginDir == "" {
		logger.Warn("Plugin directory not set, skipping plugin loading")
		return nil
	}

	// Find all .so files in the plugin directory
	pluginPaths, err := filepath.Glob(filepath.Join(pm.pluginDir, "*.so"))
	if err != nil {
		return fmt.Errorf("failed to find plugins: %w", err)
	}

	logger.Info("Found plugins", logger.F("count", len(pluginPaths)))

	// Load each plugin
	for _, pluginPath := range pluginPaths {
		if err := pm.LoadPlugin(pluginPath); err != nil {
			logger.Error("Failed to load plugin",
				logger.F("path", pluginPath),
				logger.F("error", err))
		}
	}

	return nil
}

// LoadPlugin loads a single plugin
func (pm *PluginManager) LoadPlugin(pluginPath string) error {
	// Check if plugin is already loaded
	pluginName := filepath.Base(pluginPath)

	pm.pluginsMutex.RLock()
	_, exists := pm.loadedPlugins[pluginName]
	pm.pluginsMutex.RUnlock()

	if exists {
		logger.Warn("Plugin already loaded", logger.F("name", pluginName))
		return nil
	}

	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", pluginName, err)
	}

	// Store the loaded plugin
	pm.pluginsMutex.Lock()
	pm.loadedPlugins[pluginName] = p
	pm.pluginsMutex.Unlock()

	// Look for Init symbol
	initSym, err := p.Lookup("Init")
	if err != nil {
		return fmt.Errorf("plugin %s does not export Init symbol: %w", pluginName, err)
	}

	// Call Init function
	initFunc, ok := initSym.(func(*PluginManager) error)
	if !ok {
		return fmt.Errorf("plugin %s has invalid Init function signature", pluginName)
	}

	if err := initFunc(pm); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", pluginName, err)
	}

	logger.Info("Successfully loaded plugin", logger.F("name", pluginName))
	return nil
}

// RegisterExtractor registers an extractor
func (pm *PluginManager) RegisterExtractor(name string, extractor interfaces.Extractor) {
	pm.extractorsMutex.Lock()
	defer pm.extractorsMutex.Unlock()

	pm.extractors[name] = extractor
	logger.Info("Registered extractor", logger.F("name", name))
}

// GetExtractor retrieves an extractor by name
func (pm *PluginManager) GetExtractor(name string) (interfaces.Extractor, bool) {
	pm.extractorsMutex.RLock()
	defer pm.extractorsMutex.RUnlock()

	extractor, exists := pm.extractors[name]
	return extractor, exists
}

// GetAllExtractors returns all registered extractors
func (pm *PluginManager) GetAllExtractors() map[string]interfaces.Extractor {
	pm.extractorsMutex.RLock()
	defer pm.extractorsMutex.RUnlock()

	// Create a copy to avoid concurrent map access
	extractors := make(map[string]interfaces.Extractor, len(pm.extractors))
	for name, extractor := range pm.extractors {
		extractors[name] = extractor
	}

	return extractors
}

// RegisterFormatter registers a formatter
func (pm *PluginManager) RegisterFormatter(name string, formatter interfaces.FormatterService) {
	pm.formattersMutex.Lock()
	defer pm.formattersMutex.Unlock()

	pm.formatters[name] = formatter
	logger.Info("Registered formatter", logger.F("name", name))
}

// GetFormatter retrieves a formatter by name
func (pm *PluginManager) GetFormatter(name string) (interfaces.FormatterService, bool) {
	pm.formattersMutex.RLock()
	defer pm.formattersMutex.RUnlock()

	formatter, exists := pm.formatters[name]
	return formatter, exists
}

// GetAllFormatters returns all registered formatters
func (pm *PluginManager) GetAllFormatters() map[string]interfaces.FormatterService {
	pm.formattersMutex.RLock()
	defer pm.formattersMutex.RUnlock()

	// Create a copy to avoid concurrent map access
	formatters := make(map[string]interfaces.FormatterService, len(pm.formatters))
	for name, formatter := range pm.formatters {
		formatters[name] = formatter
	}

	return formatters
}
