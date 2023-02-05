/* SuchMiner
 * Copyright 2022 duggavo
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

const default_config = `{
	"api": {
		"id": null,
		"worker-id": null
	},
	"http": {
		"enabled": false,
		"host": "127.0.0.1",
		"port": 0,
		"access-token": null,
		"restricted": true
	},
	"autosave": true,
	"background": false,
	"colors": false,
	"title": true,
	"randomx": {
		"init": -1,
		"init-avx2": -1,
		"mode": "auto",
		"1gb-pages": false,
		"rdmsr": true,
		"wrmsr": true,
		"cache_qos": false,
		"numa": true,
		"scratchpad_prefetch_mode": 1
	},
	"cpu": {
		"enabled": true,
		"huge-pages": true,
		"huge-pages-jit": false,
		"hw-aes": null,
		"priority": null,
		"memory-pool": false,
		"yield": true,
		"max-threads-hint": 100,
		"asm": true,
		"argon2-impl": null,
		"cn/0": false,
		"cn-lite/0": false
	},
	"opencl": {
		"enabled": false,
		"cache": true,
		"loader": null,
		"platform": "AMD",
		"adl": true,
		"cn/0": false,
		"cn-lite/0": false
	},
	"cuda": {
		"enabled": false,
		"loader": null,
		"nvml": true,
		"cn/0": false,
		"cn-lite/0": false
	},
	"donate-level": 5,
	"donate-over-proxy": 5,
	"log-file": null,
	"pools": [
		{
			"algo": null,
			"coin": "WOW",
			"url": "$URL",
			"user": "$ADDRESS",
			"spend-secret-key": "$SPEND",
			"pass": "x",
			"daemon": true,
			"rig-id": null,
			"nicehash": false,
			"keepalive": false,
			"enabled": true,
			"tls": false,
			"tls-fingerprint": null,
			"socks5": null,
			"self-select": null,
			"submit-to-origin": false
		}
	],
	"print-time": 20,
	"health-print-time": 60,
	"dmi": true,
	"retries": 10,
	"retry-pause": 5,
	"syslog": false,
	"tls": {
		"enabled": false,
		"protocols": null,
		"cert": null,
		"cert_key": null,
		"ciphers": null,
		"ciphersuites": null,
		"dhparam": null
	},
	"dns": {
		"ipv6": false,
		"ttl": 30
	},
	"user-agent": null,
	"verbose": 0,
	"watch": true,
	"pause-on-battery": false,
	"pause-on-active": false
}`
