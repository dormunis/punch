package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_InitConfig_Defaults(t *testing.T) {
	tempDir := t.TempDir()
	viper.AddConfigPath(tempDir)

	config, err := InitConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "sqlite3", config.Database.Engine)
	assert.Equal(t, filepath.Join(tempDir, "punch.db"), config.Database.Path)
}

func TestConfig_InitConfig_FromFile(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.toml")
	viper.AddConfigPath(tempDir)
	viper.SetConfigType("toml")

	configContent := `
    [settings]
    editor = "nano"
    `
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	config, err := InitConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "nano", config.Settings.Editor)
}

func TestConfig_InitConfig_AllowedRemoteUsed(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.toml")
	viper.AddConfigPath(tempDir)
	viper.SetConfigType("toml")

	remoteName := "origin"

	spreadsheetRemote := fmt.Sprintf(`
        [remotes.%s]
        type = "spreadsheet"
        spreadsheet_id = "1"
        sheet_name = "Sheet1"

        [remotes.%s.columns]
        id = "A"
        client = "B"
        date = "C"
        start_time = "D"
        end_time = "E"
        total_time = "F"
        note = "G"
        `, remoteName, remoteName)

	err := os.WriteFile(configFile, []byte(spreadsheetRemote), 0644)
	assert.NoError(t, err)

	config, err := InitConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, 1, len(config.Remotes))
	assert.Equal(t, "spreadsheet", config.Remotes[remoteName].Type())
}

func TestConfig_InitConfig_NotSupportedRemoteReturnsError(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.toml")
	viper.AddConfigPath(tempDir)
	viper.SetConfigType("toml")

	remoteName := "origin"
	invalidRemoteName := "shrek"

	invalidRemoteConfig := fmt.Sprintf(`
        [remotes.%s]
        type = "%s"
    `, remoteName, invalidRemoteName)

	err := os.WriteFile(configFile, []byte(invalidRemoteConfig), 0644)
	assert.NoError(t, err)

	config, err := InitConfig(tempDir)

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfig_InitConfig_AutoSyncValidatedWhenDefaultIsSet(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.toml")
	viper.AddConfigPath(tempDir)
	viper.SetConfigType("toml")

	remoteName := "origin"

	spreadsheetRemote := fmt.Sprintf(`
        [remotes.%s]
        type = "spreadsheet"
        spreadsheet_id = "1"
        sheet_name = "Sheet1"
        default_remote = "%s"
        autosync = ["start"]

        [remotes.%s.columns]
        id = "A"
        client = "B"
        date = "C"
        start_time = "D"
        end_time = "E"
        total_time = "F"
        note = "G"
        `, remoteName, remoteName, remoteName)

	err := os.WriteFile(configFile, []byte(spreadsheetRemote), 0644)
	assert.NoError(t, err)

	config, err := InitConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestConfig_InitConfig_AutoSyncNotValidatedWhenDefaultIsNotSet(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.toml")
	viper.AddConfigPath(tempDir)
	viper.SetConfigType("toml")

	remoteName := "origin"

	spreadsheetRemote := fmt.Sprintf(`
        [remotes.%s]
        type = "spreadsheet"
        spreadsheet_id = "1"
        sheet_name = "Sheet1"
        autosync = ["start"]

        [remotes.%s.columns]
        id = "A"
        client = "B"
        date = "C"
        start_time = "D"
        end_time = "E"
        total_time = "F"
        note = "G"
        `, remoteName, remoteName)

	err := os.WriteFile(configFile, []byte(spreadsheetRemote), 0644)
	assert.NoError(t, err)

	config, err := InitConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestConfig_InitConfig_AutoSyncNotValidatedWhenDefaultNotExists(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.toml")
	viper.AddConfigPath(tempDir)
	viper.SetConfigType("toml")

	spreadsheetRemote := `
        [settings]
        default_remote = "not-exists"
        `

	err := os.WriteFile(configFile, []byte(spreadsheetRemote), 0644)
	assert.NoError(t, err)

	config, err := InitConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, config)
}
