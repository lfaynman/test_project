package main

import (
    "fmt"
    "strings"
    "strconv"
    "github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
    //固定讀取上傳後的execl檔案
    xlsx, err := excelize.OpenFile("/tmp/Extensions_template.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    userDataSheetName := xlsx.GetSheetMap()[1]
    
    //計算總筆數
    rows, err := xlsx.GetRows(userDataSheetName)
    if err != nil {
        fmt.Println(err)
        return
    }
    var totalRecord int
    for index, _ := range rows {
                totalRecord = index+1
        }

        var mac,phoneType []string
        var forPhoneType = [...]string{"x3s","c58p"}
        for i:=3;i<=totalRecord; i++ {
            macaddr, err := xlsx.GetCellValue(userDataSheetName, "AZ"+strconv.Itoa(i))
            if err != nil {
        fmt.Println(err)
        }
        replacer := strings.NewReplacer(":", "", "-", "")
            mac = append(mac,strings.ToLower(replacer.Replace(macaddr)))
        phoneTypeD, err := xlsx.GetCellValue(userDataSheetName, "AY"+strconv.Itoa(i))
        if err != nil {
        fmt.Println(err)
        }
        phoneType = append(phoneType,phoneTypeD)
        }

        for i:=0; i<=totalRecord-3;i++{ 
            fmt.Println(mac[i],phoneType[i])
        }    
        fmt.Println(forPhoneType[0])

}
