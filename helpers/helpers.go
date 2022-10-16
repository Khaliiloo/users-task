package helpers

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"users-task/configs"
)

func CalculateAge(dob time.Time) int {
	now := time.Now()
	diff := now.Sub(dob)
	return int(math.Floor(diff.Hours() / (365.0 * 24.0)))
}

func ToJson(i interface{}) string {

	j, _ := json.Marshal(i)
	return string(j)
}

func ConvertMapToJson(m map[string][]string) string {
	res := map[string]string{}
	for k, v := range m {
		value := ""
		for _, val := range v {
			value += val + ", "
		}
		res[k] = value
	}
	return ToJson(res)
}

func GetFileNameToAdd(fileName string, userId int, index int) string {
	matches, err := filepath.Glob("./files/" + strconv.Itoa(userId) + "/" + fileName)
	if err != nil {
		configs.Logger.Error(err)
	}

	if len(matches) == 0 {
		return "./files/" + strconv.Itoa(userId) + "/" + fileName
	} else {
		fileArr := strings.Split(fileName, ".")
		if len(fileArr) == 1 {
			fileName += "_" + strconv.Itoa(index)

		} else if len(fileArr) == 2 {
			fileArr[len(fileArr)-2] += "." + strconv.Itoa(index)
			fileName = strings.Join(fileArr, ".")
		} else {
			fileArr[len(fileArr)-2] = strconv.Itoa(index)
			fileName = strings.Join(fileArr, ".")
		}
		index++
		return GetFileNameToAdd(fileName, userId, index)
	}
}

func CreateFolder(folderName string) {
	if _, err := os.Stat(folderName); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folderName, os.ModePerm)
		if err != nil {
			configs.Logger.Errorf(err.Error())
		}
	}
}

func ListFilesInDir(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	filesSlice := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			filesSlice = append(filesSlice, file.Name())
		}
	}
	return filesSlice
}
