/* SuchMiner
 * Copyright 2024 duggavo
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
	Pools          []Pool `json:"pools"`
	Wallet         string `json:"wallet"`
	SpendSecretKey string `json:"spend-secret-key"`
}

type Pool struct {
	URL    string `json:"url"`
	Daemon bool   `json:"daemon,omitempty"`
	TLS    bool   `json:"tls,omitempty"`
}

var cfg = &Config{
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
	pools := make([]XmrigPool, len(c.Pools))
	for i := range pools {
		pools[i] = c.Pools[i].toXmrig(c)
	}

	return XmrigConfig{
		Autosave:    false,
		DonateLevel: 5,
		CPU:         true,
		OpenCL:      false,
		CUDA:        false,
		Pools:       pools,
		Retries:     4,
		RetryPause:  2,
	}
}

func (c *Config) JSON() []byte {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}
