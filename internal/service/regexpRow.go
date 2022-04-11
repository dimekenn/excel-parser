package service

import(
	"regexp"
	"github.com/labstack/gommon/log"
	"strings"
	"strconv"
)

// func getNameReg (row []string) (int, error){
// 	var idSlice []int
// 	for i, v := range row{
// 		// if strings.Contains(v, "наименование"){
// 		// 	idSlice = append(idSlice, i)
// 		// }
// 		nameReg, nameRegErr := regexp.Compile(`наименование|Наименование`)
// 		if nameRegErr != nil{
// 			log.Errorf("failed to compile regexp: %s", nameRegErr)
// 			return 0, echo.NewHTTPError(http.StatusInternalServerError, nameRegErr)
// 		}

// 		if nameReg.MatchString(v) {
// 			idSlice = append(idSlice, i)
// 		}
// 	}
// 	var rowIdSlice []int
// 	if len(idSlice) > 1 {
// 		for _, v := range idSlice{
// 			totalReg, totalRegErr := regexp.Compile(`полное|Полное`)
// 			if totalRegErr != nil{
// 				log.Errorf("failed to compile regexp: %s", totalRegErr)
// 				return 0, echo.NewHTTPError(http.StatusInternalServerError, totalRegErr)
// 			}

// 			if totalReg.MatchString(row[v]){
// 				rowIdSlice = append(rowIdSlice, v)
// 			}
// 		} else if len(idSlice) < 1 {
			
// 		}
// 	}
	
// }

func takeVolume(name string) (string, float32, float32, float32) {
	regV, regErr := regexp.Compile(`\d+х.+х\d+`)
	if regErr != nil {
		log.Errorf("failed to regexp compile: %v", regErr)
		return name, 0, 0, 0
	}
	volume := regV.FindString(name)
	if volume != "" {
		vols := strings.Split(volume, "х")
		if len(vols) == 3 {
			length, lenErr := strconv.ParseFloat(strings.Replace(vols[0], ",", ".", 1), 32)
			if lenErr != nil {
				log.Errorf("failed to parse float: %v", lenErr)
				return name, 0, 0, 0
			}
			height, heiErr := strconv.ParseFloat(strings.Replace(vols[1], ",", ".", 1), 32)
			if heiErr != nil {
				log.Errorf("failed to parse float: %v", heiErr)
				return name, 0, 0, 0
			}
			width, widErr := strconv.ParseFloat(strings.Replace(vols[2], ",", ".", 1), 32)
			if widErr != nil {
				log.Errorf("failed to parse float: %v", widErr)
				return name, 0, 0, 0
			}
			name = strings.Replace(name, volume, "", 1)
			return name, float32(length), float32(height), float32(width)
		} else {
			return name, 0, 0, 0
		}
	}
	return name, 0, 0, 0
}