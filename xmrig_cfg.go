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

type XmrigConfig struct {
	Autosave    bool        `json:"autosave"`
	DonateLevel int         `json:"donate-level"`
	CPU         bool        `json:"cpu"`
	OpenCL      bool        `json:"opencl"`
	CUDA        bool        `json:"cuda"`
	Pools       []XmrigPool `json:"pools"`
	Retries     int         `json:"retries,omitempty"`
	RetryPause  int         `json:"retry-pause,omitempty"`
}
type XmrigPool struct {
	Algo           *string `json:"algo"`
	Coin           *string `json:"coin"`
	URL            string  `json:"url"`
	User           string  `json:"user"`
	SpendSecretKey *string `json:"spend-secret-key"`
	Pass           string  `json:"pass"`
	Daemon         bool    `json:"daemon"`
	TLS            bool    `json:"tls"`
	TLSFingerprint *string `json:"tls-fingerprint"`
}

func (c XmrigConfig) JSON() []byte {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	return b
}
