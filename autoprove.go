package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const phoneSetPath string = "/usr/local/voip/ipphone/"

var splitCharacterSupport = []string{"=", ":"}
var supportPhoneTypes []string

func main() {
	orignalTime := time.Now().Unix()
	var config, setPhone []string
	path, err := exec.Command("/bin/sh", "-c", `/bin/cat /etc/xinetd.d/tftp |grep server_args|awk -F" " '{print $4}'`).Output()
	if err != nil {
		writeToLog(err.Error())
		os.Exit(1)
	}
	tftpOutPutPath := strings.Join(strings.Fields(string(path)), "") + "/"
	command := "/bin/rm -f " + tftpOutPutPath + "*"
	exec.Command("/bin/sh", "-c", command).Output()
	exec.Command("/bin/sh", "-c", "true > /usr/local/voip/ipphone/autoprovlog.log").Output()
	f, err := excelize.OpenFile("/tmp/Extensions_data.xlsx")
	if err != nil {
		writeToLog(err.Error())
		return
	}
	phoneData, err := f.GetRows(f.GetSheetMap()[1])
	if err != nil {
		writeToLog(err.Error())
		return
	}
	isMobile := covertExcelItemToArrayItem("A")
	isUseAutoProvisioning := covertExcelItemToArrayItem("K")
	usePhoneType := covertExcelItemToArrayItem("AY")
	macAddress := covertExcelItemToArrayItem("AZ")
	useSetPhoneTemplate := covertExcelItemToArrayItem("BA")
	supportPhoneTypes = getSupportPhoneTypes(phoneSetPath)
	phoneData = makeFormateRegular(phoneData[2:], isMobile, isUseAutoProvisioning, usePhoneType, macAddress)

	/*	ch := make(chan string, len(supportPhoneType))
		for _, set := range supportPhoneType {
			go func(set string) {
				ch <- set
			}(set)
		}
		for range supportPhoneType {
			fmt.Println(<-ch)
		}
	*/

	for i, _ := range phoneData {
		if phoneData[i][usePhoneType] == "unknowtype" || phoneData[i][isUseAutoProvisioning] == "no" || phoneData[i][isMobile] == "yes" {
			continue
		}
		if i == 0 || phoneData[i][usePhoneType] != phoneData[i][usePhoneType-1] || phoneData[i][useSetPhoneTemplate] != phoneData[i][useSetPhoneTemplate-1] {
			config = config[:0]
			configlines, err := readLines(phoneSetPath + phoneData[i][usePhoneType] + "/" + phoneData[i][useSetPhoneTemplate])
			if err != nil {
				writeToLog(err.Error())
			}
			for _, cline := range configlines {
				config = append(config, cline)
			}
			setPhone = setPhone[:0]
			setPhonelines, err := readLines(phoneSetPath + phoneData[i][usePhoneType] + "/setphone")
			if err != nil {
				writeToLog(err.Error())
			}

			for _, sline := range setPhonelines {
				if strings.Join(strings.Fields(sline), "") == "" {
					continue
				}
				setPhone = append(setPhone, sline)
			}
		}
		for _, item := range setPhone {
			var configTargetItem int
			var phoneSetValue, splitCharacter string
			setPhoneSearchString := strings.Join(strings.Split(item, ":")[:1], "")
			setPhoneData := strings.Split(strings.Join(strings.Split(item, ":")[1:2], ""), ",")
			for _, data := range setPhoneData {
				phoneSetValue = phoneSetValue + strings.Replace(data, strings.Join(strings.Fields(data), ""), phoneData[i][covertExcelItemToArrayItem(strings.Join(strings.Fields(data), ""))], -1)
			}
			typeFound := false
			for i, configItem := range config {
				if strings.Contains(configItem, setPhoneSearchString) {
					typeFound = true
					configTargetItem = i
					continue
				}
			}
			if typeFound == false {
				writeToLog(phoneData[i][usePhoneType] + ": Searchstring  [ " + setPhoneSearchString + " ]  ; " + "The setPhone file has searched no mach with setphone item , check the item of setphone is correct")
			}
			for _, c := range splitCharacterSupport {
				if strings.Contains(config[len(config)/2], c) {
					splitCharacter = c
					break
				}
			}
			config[configTargetItem] = strings.Join(strings.SplitAfter(config[configTargetItem], splitCharacter)[:1], "") + phoneSetValue
		}

		phoneTypefirstCharacter := phoneData[i][usePhoneType][:1]
		fileNameExtension := "." + strings.Join(strings.Split(phoneData[i][usePhoneType], ".")[1:2], "")
		switch {
		case "A" <= phoneTypefirstCharacter && phoneTypefirstCharacter <= "Z", "0" <= phoneTypefirstCharacter && phoneTypefirstCharacter <= "9":
			if err := writeLines(config, tftpOutPutPath+strings.ToUpper(phoneData[i][macAddress])+fileNameExtension); err != nil {
				writeToLog("Writeconfigfile Error : " + err.Error())
			}
		default:
			if err := writeLines(config, tftpOutPutPath+phoneData[i][macAddress]+fileNameExtension); err != nil {
				writeToLog("Writeconfigfile Error : " + err.Error())
			}
		}
	}
	writeToLog("Total " + strconv.Itoa(len(phoneData)) + "  Execel items maked.")
	writeToLog("Total RunTime is  " + strconv.Itoa(int(time.Now().Unix()-orignalTime)) + "  sec.")
}

func covertExcelItemToArrayItem(s string) int {
	var sum int
	p := len(s) - 1
	for _, r := range strings.ToUpper(s) {
		sum = (int(r)-64)*int(math.Pow(26, float64(p))) + sum
		if p > 0 {
			p--
		}
	}
	return (sum - 1)
}

func getSupportPhoneTypes(path string) []string {
	var phoneType []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fName, err := os.Stat(path + f.Name())
		if err != nil {
			writeToLog(err.Error())
		}
		if fName.Mode().IsDir() {
			phoneType = append(phoneType, fName.Name())
		}
	}
	return phoneType
}

func makeFormateRegular(phoneData [][]string, isMobile, isUseAutoProvisioning, usePhoneType, macAddress int) [][]string {
	replacer := strings.NewReplacer(":", "", "-", "")
	for i, _ := range phoneData {
		phoneData[i][macAddress] = strings.ToLower(replacer.Replace(phoneData[i][macAddress]))
		phoneData[i][isMobile] = strings.ToLower(phoneData[i][isMobile])
		phoneData[i][isUseAutoProvisioning] = strings.ToLower(phoneData[i][isUseAutoProvisioning])
		typeFound := false
		for _, pType := range supportPhoneTypes {
			if strings.Contains(strings.ToLower(phoneData[i][usePhoneType]), strings.ToLower(strings.Join(strings.Split(pType, ".")[:1], ""))) {
				typeFound = true
				phoneData[i][usePhoneType] = pType
				continue
			}
		}
		if typeFound == false {
			writeToLog("[The execl row '" + strconv.Itoa(i+3) + "']  phone type is  '" + phoneData[i][usePhoneType] + "'  the type not yet suppert so check The Phone Type\n [" + strconv.Itoa(i+3) + "] " + strings.Join(phoneData[i], " "))
			phoneData[i][usePhoneType] = "unknowtype"
		}
	}
	return phoneData
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func writeToLog(logmessage string) {
	f, err := os.OpenFile("/usr/local/voip/ipphone/autoprovlog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(logmessage)
}
