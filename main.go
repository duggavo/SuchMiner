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

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var debug = false

var Log *log.Logger

const VERSION = "2.0.0"

var hashrate float64 = 0
var difficulty float64 = 0
var poolUrl string = ""

var hashrateRegex = regexp.MustCompile(`10s\/60s\/15m [\d\.n/a]+ [\d\.n/a]+`)

var addressRegex = regexp.MustCompile("^Wo[0-9AB][123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{94}$")
var spendKeyRegex = regexp.MustCompile("^[0-9a-fA-F]{64}$")

const detect_donate = "wowrig.mooo.com"
const deps_folder = "deps_" + runtime.GOOS

var logreader = LogReader{}

func main() {
	saveConfigAndQuit := false

	flag.BoolVar(&debug, "debug", false, "runs SuchMiner in debug mode")
	flag.BoolVar(&cfg.NonInteractive, "non-interactive", false, "runs SuchMiner in non-interactive mode")
	flag.StringVar(&cfg.LogFile, "log-file", cfg.LogFile, "defines the log file")
	flag.StringVar(&cfg.Wallet, "wallet", cfg.Wallet, "wownero wallet address to mine at")
	flag.StringVar(&cfg.SpendSecretKey, "spend-secret-key", cfg.SpendSecretKey, "secret spend key")
	var poolsJson string
	flag.StringVar(&poolsJson, "pools", "", "the pools to use in JSON array format")

	flag.BoolVar(&saveConfigAndQuit, "save-config-and-quit", false, "Saves the configuration provided by command-line flags and quits")
	flag.Parse()

	loadCfgErr := loadConfig()

	if poolsJson != "" {
		err := json.Unmarshal([]byte(poolsJson), &cfg.Pools)
		if err != nil {
			fmt.Println(Red, err, Reset)
		}
	}

	if saveConfigAndQuit {
		saveConfig()
		fmt.Println("configuration file saved, quitting")
		return
	}

	if cfg.LogFile != "" {
		var err error
		logreader.File, err = os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Println(Red, err, Reset)
		} else {
			fmt.Println("Using log file:", cfg.LogFile)
		}
	}

	if debug {
		Log = log.New(logreader, Reset, log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Log = log.New(logreader, Reset, log.Ldate|log.Ltime)
	}

	Log.Println(Cyan+"Starting", Magenta+Bright+"Such"+Yellow+"Miner "+Reset+Cyan+"v"+VERSION+Reset)

	if loadCfgErr != nil {
		if debug {
			Log.Println(Red, loadCfgErr, Reset)
		}
		if cfg.Wallet == "" || len(cfg.Pools) == 0 || cfg.SpendSecretKey == "" {
			if cfg.NonInteractive {
				Log.Println(Red, "SuchMiner is not correctly configured and non-interactive mode is enabled"+Reset)
				os.Exit(1)
				return
			} else {
				configPrompt()
			}
		}
		if err := saveConfig(); err != nil {
			panic(err)
		}
	} else {
		if err := saveConfig(); err != nil {
			Log.Println(Red + err.Error() + Reset)
		}
	}

	var wowrig *exec.Cmd
	if runtime.GOOS == "windows" {
		wowrig = exec.Command(deps_folder+"/xmrig.exe", "--no-color")
	} else {
		wowrig = exec.Command(deps_folder+"/xmrig", "--no-color")
	}

	var outb saveOutput
	wowrig.Stdout = &outb

	Log.Println(Cyan + "Miner started")

	err := wowrig.Run()
	if err != nil {
		Log.Println(Bright+Red+"error running wowrig:", err)
	}
}

func configPrompt() {
	for {
		cfg.Wallet = Prompt(Cyan + "Wownero primary address (starts with Wo...): " + Reset)
		if !addressRegex.Match([]byte(cfg.Wallet)) {
			Log.Println(Red + "Error: address '" + cfg.Wallet + "' is not valid" + Reset)
			continue
		}
		break
	}
	for {
		cfg.SpendSecretKey = Prompt(Cyan + "Secret spend key: " + Reset)
		if !spendKeyRegex.Match([]byte(cfg.SpendSecretKey)) {
			Log.Println(Bright + Red + "Error: spend key is not valid" + Reset)
			continue
		}
		break
	}
	for {
		log.Println("SuchMiner has a default node list, which tries to connect to the local daemon first," +
			"or to a public node in case of failure." + Reset)
		isRunningDaemon := Prompt(Cyan + "Do you want to use the default node list? (y/N) " + Reset)
		if !strings.HasPrefix(strings.ToLower(isRunningDaemon), "y") {
			kind := Prompt(Cyan + "are you trying to connect to a daemon or proxy? (0: daemon, 1: proxy) " + Reset)
			if kind != "0" && kind != "1" {
				Log.Println(Red + "invalid choice, enter 0 or 1" + Reset)
				continue
			}

			daemonUrl := Prompt(Cyan + "Enter the daemon/proxy's IP:PORT: " + Reset)
			if len(daemonUrl) < 2 {
				Log.Println(Red + "invalid URL" + Reset)
			}

			cfg.Pools = []Pool{
				{
					URL:    daemonUrl,
					Daemon: kind == "0",
				},
			}
		}
		break
	}
}

func loadConfig() error {
	b, err := os.ReadFile("./such_config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, cfg)
	return err
}

func saveConfig() error {
	err := os.WriteFile("./such_config.json", cfg.JSON(), 0o666)
	if err != nil {
		return err
	}
	return os.WriteFile(deps_folder+"/config.json", cfg.ToXMRIG().JSON(), 0o666)
}

type saveOutput struct {
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	pStr := string(p)
	parseHashrate(pStr)
	parseMsr(pStr)
	parseErrors(pStr)
	parseJob(pStr)

	return logreader.File.Write(p)
}

func Prompt(l string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, l)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func parseHashrate(txt string) {
	txt = strings.Split(txt, "max")[0]

	indexes := hashrateRegex.FindIndex([]byte(txt))
	if indexes == nil {
		return
	}
	hashrateStr := strings.Split(string(txt[indexes[0]+12:indexes[0]+20]), " ")[0]
	if hashrateStr == "n/a" {
		return
	}
	hashrateNum, err := strconv.ParseFloat(hashrateStr, 64)
	if err != nil {
		panic(err)
	}

	if strings.Contains(txt, "KH/s") {
		hashrateNum *= 1000
	} else if strings.Contains(txt, "MH/s") {
		hashrateNum *= 1000 * 1000
	} else if strings.Contains(txt, "GH/s") {
		hashrateNum *= 1000 * 1000 * 1000
	}

	hashrate = hashrateNum
	Log.Println(Cyan+"Hashrate:"+Yellow, uint64(hashrate), Cyan+"H/s"+Reset)
}

func parseMsr(txt string) {
	if strings.Contains(txt, "FAILED TO APPLY MSR") {
		Log.Println(Yellow + "Failed to apply MSR mod: hashrate will be low." + Reset)
		Log.Println(Yellow + "Try running this miner as root/administrator to mine more efficiently." + Reset)
	} else if strings.Contains(txt, "preset have been set") {
		Log.Println(Green + "MSR mod has been applied! You are mining at best efficiency." + Reset)
	}
}

func parseErrors(txt string) {
	if !strings.Contains(txt, "error:") {
		return
	}

	txt = strings.Split(txt, "error:")[1]

	Log.Println(Red+"WOWRig error:", txt+Reset)
}

func parseJob(txt string) {
	if strings.Contains(txt, "accepted") {
		if !strings.Contains(poolUrl, detect_donate) {
			Log.Println(Bright + Green + "BLOCK FOUND!")
			Log.Println(Bright + Green + strings.Split(txt, "accepted")[1] + Reset)
		}
		return
	}

	if !strings.Contains(txt, "new job from") {
		return
	}
	s := strings.Split(strings.Split(txt, "new job from ")[1], " ")

	nodeUrl := s[0]
	rawDiff := s[2]
	algo := s[4]

	poolUrl = nodeUrl

	diff, err := parseDiff(strings.ToLower(rawDiff))
	if err != nil {
		Log.Println(Red, err, Reset)
	}

	Log.Println(Cyan+"New job from", nodeUrl, "diff", rawDiff, "algo", algo+Reset)

	if hashrate != 0 && difficulty != diff && diff != 0 {
		difficulty = diff

		timeSecs := difficulty / hashrate
		blockTime := timeSecs / 60 / 60
		const base = "With the current difficulty and hashrate, you are expected to find a block every"
		if blockTime < 24*7 {
			Log.Println(base+Yellow, math.Round(blockTime), "hours"+Reset, "in average")
		} else {
			Log.Println(base+Yellow, math.Round(blockTime/24), "days"+Reset, "in average")
		}
	}

	difficulty = diff
}

func parseDiff(d string) (float64, error) {
	var mult float64 = 0
	if strings.HasSuffix(d, "g") {
		d = strings.TrimSuffix(d, "g")
		mult = 1_000_000_000
	}
	if strings.HasSuffix(d, "m") {
		d = strings.TrimSuffix(d, "m")
		mult = 1_000_000
	}
	if strings.HasSuffix(d, "k") {
		d = strings.TrimSuffix(d, "k")
		mult = 1_000
	}

	v, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0, err
	}
	return v * mult, err
}
