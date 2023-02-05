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

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var debug = false

var Log *log.Logger

const VERSION = "1.0.0"

const Reset = "\u001b[0m"
const Red = "\u001b[31m"
const Green = "\u001b[32m"
const Yellow = "\u001b[33m"
const Magenta = "\u001b[35m"
const Cyan = "\u001b[36m"
const Bright = "\u001b[1m"

var hashrate float64 = 0

var hashrateRegex = regexp.MustCompile(`10s\/60s\/15m [\d\.n/a]+ [\d\.n/a]+`)

var addressRegex = regexp.MustCompile("^W[Wo][0-9AB][123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{94}$")
var spendKeyRegex = regexp.MustCompile("^[0-9a-fA-F]{64}$")

var address string
var spendKey string
var url string

func main() {
	flag.BoolVar(&debug, "debug", false, "")
	flag.Parse()

	if debug {
		Log = log.New(os.Stdout, Reset, log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Log = log.New(os.Stdout, Reset, log.Ldate|log.Ltime)
	}

	Log.Println(Reset + Cyan)
	Log.Println(Cyan+"Starting", Magenta+Bright+"Such"+Yellow+"Miner "+Reset+Cyan+"v"+VERSION)
	Log.Println(Reset)

	if _, err := os.Stat("./deps/config.json"); errors.Is(err, os.ErrNotExist) {
		getConfig()
		setConfig()
	}

	var wowrig *exec.Cmd
	if runtime.GOOS == "windows" {
		wowrig = exec.Command("./deps/wowrig.exe", "--no-color")
	} else {
		wowrig = exec.Command("./deps/wowrig", "--no-color")
	}

	var outb saveOutput
	wowrig.Stdout = &outb

	Log.Println(Cyan + "Miner started")

	wowrig.Run()
}

func getConfig() {

	address = Prompt(Cyan + "Wownero address: " + Reset)
	spendKey = Prompt(Cyan + "Spend key: " + Reset)
	isRunningDaemon := Prompt(Cyan + "Are you running a Wownero daemon? (y/n) " + Reset)
	if strings.HasPrefix(isRunningDaemon, "y") {
		url = "127.0.0.1:34568"
	} else {
		url = "wowrig.mooo.com:34568"
	}

	if !addressRegex.Match([]byte(address)) {
		fmt.Println(Red + "Error: address '" + address + "' is not valid")
		return
	} else if !spendKeyRegex.Match([]byte(spendKey)) {
		fmt.Println(Red + "Error: spend key is not valid")
		return
	}

}

func setConfig() {
	config := default_config

	config = strings.Replace(config, "$ADDRESS", address, 1)
	config = strings.Replace(config, "$SPEND", spendKey, 1)
	config = strings.Replace(config, "$URL", url, 1)

	os.WriteFile("./deps/config.json", []byte(config), 0600)
}

type saveOutput struct {
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	parseHashrate(string(p))

	parseMsr(string(p))

	if debug {
		os.Stdout.Write(p)
	}
	return len(p), nil
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
	Log.Println(Reset+Cyan+"Hashrate:"+Yellow, uint64(hashrate), Cyan+"H/s")
}

func parseMsr(txt string) {
	if strings.Contains(txt, "FAILED TO APPLY MSR") {
		Log.Println(Reset + Red + "Failed to apply MSR mod: hashrate will be low.")
		Log.Println(Red + "Try running this miner as root." + Reset)
	} else if strings.Contains(txt, "preset have been set") {
		Log.Println(Reset + Green + "MSR mod has been applied! You are mining at best efficiency.")
	}
}
