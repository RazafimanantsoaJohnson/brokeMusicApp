package logging

import (
	"fmt"
	"os"
	"time"
)

func LogData(newLog string) error {
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	currentLogFileName := fmt.Sprintf("%v%v%v.log", year, month, day)
	_, err := os.Stat(currentLogFileName)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(currentLogFileName)
			if err != nil {
				fmt.Printf("unable to create file named: %v\n", currentLogFileName)
			}
		}
	}
	logString := fmt.Sprintf("%v===%v\n", currentTime, newLog)
	logFile, err := os.Open(currentLogFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer logFile.Close()
	currentFileContent, err := os.ReadFile(currentLogFileName)
	if err != nil {
		return err
	}
	currentFileContent = append(currentFileContent, []byte(logString)...)
	err = os.WriteFile(currentLogFileName, currentFileContent, 0666)
	if err != nil {
		fmt.Println(logString)
		return err
	}
	return fmt.Errorf("%v, %v", string(currentFileContent), logFile)
}
