package config

import (
	"io"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Viper is a prioritized configuration registry. It
// maintains a set of configuration sources, fetches
// values to populate those, and provides them according
// to the source's priority.
// The priority of the sources is the following:
// 1. overrides
// 2. flags
// 3. env. variables
// 4. config file
// 5. key/value store
// 6. defaults
type Viper = viper.Viper

// FlagValue is an interface that users can implement
// to bind different flags to viper.
type FlagValue = viper.FlagValue

// FlagValueSet is an interface that users can implement
// to bind a set of flags to viper.
type FlagValueSet = viper.FlagValueSet

// config is default config
var config = viper.New()

func init() {
	SetConfigName("config")
	SetConfigType("yaml")
	AddConfigPath("./")
	SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
}

// OnConfigChange set handler
func OnConfigChange(run func(in fsnotify.Event)) { config.OnConfigChange(run) }

// WatchConfig for changes
func WatchConfig() { config.WatchConfig() }

// SetConfigFile explicitly defines the path, name and extension of the config file.
// Viper will use this and not check any of the config paths.
func SetConfigFile(in string) { config.SetConfigFile(in) }

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
// E.g. if your prefix is "spf", the env registry will look for env
// variables that start with "SPF_".
func SetEnvPrefix(in string) { config.SetEnvPrefix(in) }

// ConfigFileUsed returns the file used to populate the config registry.
func ConfigFileUsed() string { return config.ConfigFileUsed() }

// AddConfigPath adds a path for Viper to search for the config file in.
// Can be called multiple times to define multiple search paths.
func AddConfigPath(in string) { config.AddConfigPath(in) }

// AddRemoteProvider adds a remote configuration source.
// Remote Providers are searched in the order they are added.
// provider is a string value, "etcd" or "consul" are currently supported.
// endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
// path is the path in the k/v store to retrieve configuration
// To retrieve a config file called myapp.json from /configs/myapp.json
// you should set path to /configs and set config name (SetConfigName()) to
// "myapp"
func AddRemoteProvider(provider, endpoint, path string) error {
	return config.AddRemoteProvider(provider, endpoint, path)
}

// AddSecureRemoteProvider adds a remote configuration source.
// Secure Remote Providers are searched in the order they are added.
// provider is a string value, "etcd" or "consul" are currently supported.
// endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
// secretkeyring is the filepath to your openpgp secret keyring.  e.g. /etc/secrets/myring.gpg
// path is the path in the k/v store to retrieve configuration
// To retrieve a config file called myapp.json from /configs/myapp.json
// you should set path to /configs and set config name (SetConfigName()) to
// "myapp"
// Secure Remote Providers are implemented with github.com/xordataexchange/crypt
func AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	return config.AddSecureRemoteProvider(provider, endpoint, path, secretkeyring)
}

// SetTypeByDefaultValue enables or disables the inference of a key value's
// type when the Get function is used based upon a key's default value as
// opposed to the value returned based on the normal fetch logic.
//
// For example, if a key has a default value of []string{} and the same key
// is set via an environment variable to "a b c", a call to the Get function
// would return a string slice for the key if the key's type is inferred by
// the default value and the Get function would return:
//
//   []string {"a", "b", "c"}
//
// Otherwise the Get function would return:
//
//   "a b c"
func SetTypeByDefaultValue(enable bool) { config.SetTypeByDefaultValue(enable) }

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} { return config.Get(key) }

// Sub returns new Viper instance representing a sub tree of this instance.
// Sub is case-insensitive for a key.
func Sub(key string) *Viper { return config.Sub(key) }

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return config.GetString(key) }

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool { return config.GetBool(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return config.GetInt(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return config.GetInt64(key) }

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 { return config.GetFloat64(key) }

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return config.GetTime(key) }

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration { return config.GetDuration(key) }

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string { return config.GetStringSlice(key) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} { return config.GetStringMap(key) }

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string { return config.GetStringMapString(key) }

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return config.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(key string) uint { return config.GetSizeInBytes(key) }

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func UnmarshalKey(key string, rawVal interface{}) error { return config.UnmarshalKey(key, rawVal) }

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func Unmarshal(rawVal interface{}) error { return config.Unmarshal(rawVal) }

// UnmarshalExact unmarshals the config into a Struct, erroring if a field is nonexistent
// in the destination struct.
func UnmarshalExact(rawVal interface{}) error { return config.UnmarshalExact(rawVal) }

// BindPFlags binds a full flag set to the configuration, using each flag's long
// name as the config key.
func BindPFlags(flags *pflag.FlagSet) error { return config.BindPFlags(flags) }

// BindPFlag binds a specific key to a pflag (as used by cobra).
// Example (where serverCmd is a Cobra instance):
//
//	 serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
//	 Viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
//
func BindPFlag(key string, flag *pflag.Flag) error { return config.BindPFlag(key, flag) }

// BindFlagValues binds a full FlagValue set to the configuration, using each flag's long
// name as the config key.
func BindFlagValues(flags FlagValueSet) (err error) { return config.BindFlagValues(flags) }

// BindFlagValue binds a specific key to a FlagValue.
// Example (where serverCmd is a Cobra instance):
//
//	 serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
//	 Viper.BindFlagValue("port", serverCmd.Flags().Lookup("port"))
//
func BindFlagValue(key string, flag FlagValue) error { return config.BindFlagValue(key, flag) }

// BindEnv binds a Viper key to a ENV variable.
// ENV variables are case sensitive.
// If only a key is provided, it will use the env key matching the key, uppercased.
// EnvPrefix will be used when set when env name is not provided.
func BindEnv(input ...string) error { return config.BindEnv(input...) }

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
func IsSet(key string) bool { return config.IsSet(key) }

// AutomaticEnv has Viper check ENV variables for all.
// keys set in config, default & flags
func AutomaticEnv() { config.AutomaticEnv() }

// SetEnvKeyReplacer sets the strings.Replacer on the viper object
// Useful for mapping an environmental variable to a key that does
// not match it.
func SetEnvKeyReplacer(r *strings.Replacer) { config.SetEnvKeyReplacer(r) }

// RegisterAlias provide another accessor for the same key.
// This enables one to change a name without breaking the application
func RegisterAlias(alias string, key string) { config.RegisterAlias(alias, key) }

// InConfig checks to see if the given key (or an alias) is in the config file.
func InConfig(key string) bool { return config.InConfig(key) }

// SetDefault sets the default value for this key.
// SetDefault is case-insensitive for a key.
// Default only used when no value is provided by the user via flag, config or ENV.
func SetDefault(key string, value interface{}) { config.SetDefault(key, value) }

// Set sets the value for the key in the override regiser.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
func Set(key string, value interface{}) { config.Set(key, value) }

// ReadInConfig will discover and load the configuration file from disk
// and key/value stores, searching in one of the defined paths.
func ReadInConfig() error { return config.ReadInConfig() }

// MergeInConfig merges a new configuration with an existing config.
func MergeInConfig() error { return config.MergeInConfig() }

// ReadConfig will read a configuration file, setting existing keys to nil if the
// key does not exist in the file.
func ReadConfig(in io.Reader) error { return config.ReadConfig(in) }

// MergeConfig merges a new configuration with an existing config.
func MergeConfig(in io.Reader) error { return config.MergeConfig(in) }

// WriteConfig writes the current configuration to a file.
func WriteConfig() error { return config.WriteConfig() }

// SafeWriteConfig writes current configuration to file only if the file does not exist.
func SafeWriteConfig() error { return config.SafeWriteConfig() }

// WriteConfigAs writes current configuration to a given filename.
func WriteConfigAs(filename string) error { return config.WriteConfigAs(filename) }

// SafeWriteConfigAs writes current configuration to a given filename if it does not exist.
func SafeWriteConfigAs(filename string) error { return config.SafeWriteConfigAs(filename) }

// ReadRemoteConfig attempts to get configuration from a remote source
// and read it in the remote configuration registry.
func ReadRemoteConfig() error { return config.ReadRemoteConfig() }

// WatchRemoteConfig changes
func WatchRemoteConfig() error { return config.WatchRemoteConfig() }

// WatchRemoteConfigOnChannel changes
func WatchRemoteConfigOnChannel() error { return config.WatchRemoteConfigOnChannel() }

// AllKeys returns all keys holding a value, regardless of where they are set.
// Nested keys are returned with a v.keyDelim (= ".") separator
func AllKeys() []string { return config.AllKeys() }

// AllSettings merges all settings and returns them as a map[string]interface{}.
func AllSettings() map[string]interface{} { return config.AllSettings() }

// SetFs sets the filesystem to use to read configuration.
func SetFs(fs afero.Fs) { config.SetFs(fs) }

// SetConfigName sets name for the config file.
// Does not include extension.
func SetConfigName(in string) { config.SetConfigName(in) }

// SetConfigType sets the type of the configuration returned by the
// remote source, e.g. "json".
func SetConfigType(in string) { config.SetConfigType(in) }

// Debug prints all configuration registries for debugging
// purposes.
func Debug() { config.Debug() }
