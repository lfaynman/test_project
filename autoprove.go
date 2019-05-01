package main

import (
    "fmt"
    "strings"
    "strconv"
    "os"
    "bufio"
    "github.com/360EntSecGroup-Skylar/excelize"
)

const phoneSetPath string = "/usr/local/voip/iphone"
var forPhoneType = [...]string{"x3s","c58p","168ge"}

func ReadExcelData(column string) []string{
    
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
        totalRecord = index+1
    }
    for i:=3;i<=totalRecord; i++ {
        excelData, err := xlsx.GetCellValue(sheetName,column+strconv.Itoa(i))
        if err != nil {
        fmt.Println(err)
        }
        
        switch column {
            case "AZ" : {
                replacer := strings.NewReplacer(":", "", "-", "")
                data = append(data,strings.ToLower(replacer.Replace(excelData)))
            }
            case "A","K","BA" :{
                data = append(data,strings.ToLower(excelData))
            }
            case "AY" :{
                for _, pType := range forPhoneType {
                    if strings.Contains(strings.ToLower(excelData), pType){
                        data = append(data,pType)
                    }
                }
            }
            default : 
                data = append(data,excelData)
        }
    }
    return data
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

func main() {
    
    isMobile := ReadExcelData("A")
    phoneNum := ReadExcelData("B")
    engName := ReadExcelData("C")
    zhName := ReadExcelData("D")
    dpartName := ReadExcelData("G")
    passwd := ReadExcelData("I")
    autoprov := ReadExcelData("K")
    phoneType := ReadExcelData("AY")
    mac := ReadExcelData("AZ")
    pSetTemplate := ReadExcelData("BA")
    pSetBindIp := ReadExcelData("BB")
    
    x3sDefault, err := readLines(phoneSetPath+"/x3s/default") //check phoneSetTemplate and read the config
    if err != nil {
    fmt.Printf("readLines: %s", err)
    }
    
    var presencePhoneType []string
    for i:=0; i<len(forPhoneType);i++ {
        for j:=0; j<len(pSetTemplate); j++ {
            if forPhoneType[i] == phoneType[j] {
                presencePhoneType = append(presencePhoneType,phoneType[j])
                break
            }
        }
    }
    
    
    for i:=0; i<len(phoneNum);i++ {
    
    fmt.Println(isMobile[i],phoneNum[i],engName[i],zhName[i],dpartName[i],passwd[i],autoprov[i],mac[i],phoneType[i],pSetTemplate[i],pSetBindIp[i])
    
        
    }
    fmt.Println(presencePhoneType)
    fmt.Println(x3sDefault[6])
    
     

}
