/*
Copyright 2013 Petru Ciobanu, Francesco Paglia, Lorenzo Pierfederici

This file is part of maponet/utils.

maponet/utils is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

maponet/utils is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with maponet/utils.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
Package conf contains a simple loader for configuration files.
*/
package conf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"
)

type NotFoundError struct {
	key string
}

func (e *NotFoundError) Error() string { return fmt.Sprintf("Key not found: %s", e.key) }

type ConfigSyntaxError struct {
	line string
}

func (e *ConfigSyntaxError) Error() string {
	return fmt.Sprintf("Config syntax error on line: %s", e.line)
}

// Regexp patterns
var (
	PatternOption, _  = regexp.Compile("(.*)=([^#]*)")
	PatternComment, _ = regexp.Compile("^#")
	PatternEmpty, _   = regexp.Compile("^[\t\n\f\r ]$")
	PatternSection, _ = regexp.Compile("^\\[(.*)\\]([^#]*)")
)

// Config is a simple synchronized object that can be used to parse
// configuration files and store key=value pairs
type Config struct {
	sections map[string]map[string]string
	mutex    sync.RWMutex
}

func NewConfig() *Config {
	var c Config
	c.sections = make(map[string]map[string]string)
	return &c
}

// Parse
func (c *Config) Parse(r *bufio.Reader) error {
	var section string = "default"
	for {
		l, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		switch {
		case PatternEmpty.Match(l):
			continue
		case PatternComment.Match(l):
			continue
		case PatternSection.Match(l):
			m := PatternSection.FindSubmatch(l)
			section = string(bytes.TrimSpace(m[1]))
		case PatternOption.Match(l):
			m := PatternOption.FindSubmatch(l)
			c.Set(section, string(bytes.TrimSpace(m[1])), string(bytes.TrimSpace(m[2])))
		default:
			return &ConfigSyntaxError{string(l)}
		}
	}
	return nil
}

// ParseFile
func (c *Config) ParseFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = c.Parse(bufio.NewReader(f))
	if err != nil {
		return err
	}
	return nil
}

// Set
func (c *Config) Set(section, key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.sections[section]; !ok {
		c.sections[section] = make(map[string]string)
	}

	c.sections[section][key] = value
}

// GetString
func (c *Config) GetString(section, key string) (string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if _, ok := c.sections[section]; !ok {
		return "", &NotFoundError{section}
	}

	v, ok := c.sections[section][key]
	if !ok {
		return v, &NotFoundError{key}
	}
	return v, nil
}

// GetBool
func (c *Config) GetBool(section, key string) (bool, error) {
	s, err := c.GetString(section, key)
	if err != nil {
		return *new(bool), err
	}

	return strconv.ParseBool(s)
}

// GetInt
func (c *Config) GetInt(section, key string) (int, error) {
	s, err := c.GetString(section, key)
	if err != nil {
		return *new(int), err
	}

	return strconv.Atoi(s)
}

// GetFloat64
func (c *Config) GetFloat64(section, key string) (float64, error) {
	s, err := c.GetString(section, key)
	if err != nil {
		return *new(float64), err
	}

	return strconv.ParseFloat(s, 64)
}

var DefaultConfig = NewConfig()

func ParseFile(path string) error {
	return DefaultConfig.ParseFile(path)
}

func Set(section, key, value string) {
	DefaultConfig.Set(section, key, value)
}

func GetString(section, key string) (string, error) {
	return DefaultConfig.GetString(section, key)
}

func GetBool(section, key string) (bool, error) {
	return DefaultConfig.GetBool(section, key)
}

func GetInt(section, key string) (int, error) {
	return DefaultConfig.GetInt(section, key)
}

func GetFloat64(section, key string) (float64, error) {
	return DefaultConfig.GetFloat64(section, key)
}
