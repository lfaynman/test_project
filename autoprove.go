package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"strconv"
	"strings"
)

const phoneSetPath string = "/usr/local/voip/iphone/"

var phoneType = []string{"x3s", "c58p", "168ge"}


func ReadExcelData(column string) []string {

	var data []string

	//固定讀取上傳後的execl檔案
	xlsx, err := excelize.OpenFile("/tmp/Extensions_template.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	sheetName := xlsx.GetSheetMap()[1]

	//計算總筆數
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var totalRecord int
	for index, _ := range rows {
		totalRecord = index + 1
	}

	for i := 3; i <= totalRecord; i++ {
		excelData, err := xlsx.GetCellValue(sheetName, column+strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
		}

		switch column {
		case "AZ":
			{
				replacer := strings.NewReplacer(":", "", "-", "")
				data = append(data, strings.ToLower(replacer.Replace(excelData)))
			}
		case "A", "K", "BA":
			{
				data = append(data, strings.ToLower(excelData))
			}
		case "AY":
			{
				typeFound := false
				for _, pType := range phoneType {
					if strings.Contains(strings.ToLower(excelData), pType) {
						typeFound = true
						data = append(data, pType)
					}
				}
				if typeFound == false {
					fmt.Println("Check The Phone Type Had not provider")
					data = append(data, "notype")
				}
			}
		default:
			data = append(data, excelData)
		}
	}
	return data
}

//func ReadExcelCellData() string {

//}

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

func main() {

	xlsIsMobile := ReadExcelData("A")
	xlsAutoprov := ReadExcelData("K")
	xlsPhoneType := ReadExcelData("AY")
	xlsmac := ReadExcelData("AZ")
	xlsSetPhoneTemplate := ReadExcelData("BA")
	xlsSetPhoneBindIp := ReadExcelData("BB")

	//        for i := 0; i < len(xlsIsMobile); i++

	for i, isMobile := range xlsIsMobile {
		if isMobile == "no" && xlsAutoprov[i] == "yes" && xlsPhoneType[i] != "notype" {
			fmt.Println(i, isMobile, xlsAutoprov[i], xlsmac[i], xlsPhoneType[i], xlsSetPhoneTemplate[i], xlsSetPhoneBindIp[i])


			if i == 0 || !(xlsPhoneType[i] == xlsPhoneType[i-1] && xlsSetPhoneTemplate[i] == xlsSetPhoneTemplate[i-1]) && xlsAutoprov[i-1] == "yes" && xlsPhoneType[i-1] != "notype" {

				config, err := readLines(phoneSetPath + xlsPhoneType[i] + "/" + xlsSetPhoneTemplate[i])
				if err != nil {
					fmt.Printf("readLines: %s", err)
				}
				fmt.Printf("%s\n type= %T",config,config)
			}

			if i == 0 || !(xlsPhoneType[i] == xlsPhoneType[i-1]) && xlsAutoprov[i-1] == "yes" && xlsPhoneType[i-1] != "notype" {

				setPhone, err := readLines(phoneSetPath + xlsPhoneType[i] + "/setphone")
				if err != nil {
					fmt.Printf("readLines: %s", err)
				}

				fmt.Println(setPhone[6])
			}
		}

	}


}
