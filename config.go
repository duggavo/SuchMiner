/* SuchMiner
 * Copyright 2025 duggavo
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import "encoding/json"

type Config struct {
	Pools          []Pool        `json:"pools"`
	Wallet         string        `json:"wallet"`
	SpendSecretKey string        `json:"spend-secret-key"`
	LogFile        string        `json:"log-file,omitempty"`
	NonInteractive bool          `json:"non-interactive,omitempty"`
	HttpApi        HttpApiConfig `json:"http-api"`
}

type Pool struct {
	URL    string `json:"url"`
	Daemon bool   `json:"daemon,omitempty"`
	TLS    bool   `json:"tls,omitempty"`
}

type HttpApiConfig struct {
	Enabled     bool    `json:"enabled"`
	Host        string  `json:"host"`
	Port        uint16  `json:"port"`
	AccessToken *string `json:"access-token"`
}

var cfg = &Config{
	HttpApi: HttpApiConfig{
		Enabled:     false,
		Host:        "127.0.0.1",
		Port:        0,
		AccessToken: nil,
	},
	Pools: []Pool{
		{ // local node
			URL:    "127.0.0.1:34568",
			Daemon: true,
		},
		{ // node operated by lza_menace
			URL:    "node.suchwow.xyz:34568",
			Daemon: true,
		},
		{ // node operated by Stack Wallet
			URL:    "wownero.stackwallet.com:34568",
			Daemon: true,
		},
	},
}

func (p Pool) toXmrig(config *Config) XmrigPool {
	coin := "WOW"
	return XmrigPool{
		Coin:           &coin,
		URL:            p.URL,
		User:           config.Wallet,
		SpendSecretKey: &config.SpendSecretKey,
		Pass:           "x",
		Daemon:         p.Daemon,
		TLS:            p.TLS,
	}
}

func (c *Config) ToXMRIG() XmrigConfig {
	conf := XmrigConfig{
		Autosave:    false,
		DonateLevel: 5,
		CPU:         true,
		OpenCL:      false,
		CUDA:        false,
		Pools:       make([]XmrigPool, len(c.Pools)),
		Retries:     4,
		RetryPause:  2,
		Http:        c.HttpApi,
	}

	for i := range c.Pools {
		conf.Pools[i] = c.Pools[i].toXmrig(c)
	}

	return conf
}

func (c *Config) JSON() []byte {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}
