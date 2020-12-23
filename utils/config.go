package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Configuration is the list of GPIO ports
type Configuration struct {
	Button    Buttons
	Countdown Countdown
	Sensor    Sensors
	Laps      int
}

// Buttons are GPIO ports for butons
type Buttons struct {
	Start int
	Stop  int
}

// Countdown are GPIO ports for the countdown
type Countdown struct {
	Latch int
	Clk   int
	Data  int
	Beep  int
}

// Sensors are GPIO ports for the light sensors
type Sensors struct {
	One int
	Two int
}

func parsePath(path string) (p, n, t string) {
	p = filepath.Dir(path)
	t = filepath.Ext(path)
	n = strings.TrimSuffix(filepath.Base(path), t)
	t = strings.TrimPrefix(t, ".")
	return
}

// Get gets the configuration from file
func (c *Configuration) Get(path string) error {
	p, n, t := parsePath(path)

	// Set the file name of the configurations file
	viper.SetConfigName(n)

	// Set the path to look for the configurations file
	viper.AddConfigPath(p)

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType(t)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("Unable to decode into struct, %v", err)
	}
	return nil
}
