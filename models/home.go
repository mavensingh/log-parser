package models

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

// error levels
const (
	CRITICAL = 1
	INFO     = 2
	WARNING  = 3
)

// Logs handle the structure of the log files
type Logs struct {
	IP             string `json:"ip"`
	Timestamp      string `json:"timestamp"`
	Method         string `json:"method"`
	URL            string `json:"url"`
	ProtocolMethod string `json:"protocol_method"`
	ServerResponse string `json:"http_status"`
	SendBytes      string `json:"sent_bytes"`
	UserAgent      string `json:"user_agent"`
	// browser        string
	// system         string
}

// ErrorHandling handles the error and return formated error
func ErrorHandling(err error, msg string, code int) error {
	if err != nil {
		log.Printf("%d:%s", code, msg)
		return err
	}
	return nil
}

// OpenFile func opens a file and return back the file info as *os.File
func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ReadFile reads the file
func ReadFile(outcome *os.File) []Logs {
	defer func() {
		err := outcome.Close()
		if err != nil {
			log.Printf("error to open the file %s", err.Error())
			return
		}
	}()
	rows := []Logs{}
	scanner := bufio.NewScanner(outcome)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		matched := regularExpression(`^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+) (\S+)\s*(\S*)" (\d{3}) (\S+) (\S+) "(.*)"`, scanner.Text())
		if len(matched) == 0 {
			continue
		}
		rows = append(rows, Logs{
			IP:             matched[0][1],
			Timestamp:      matched[0][4],
			Method:         matched[0][5],
			URL:            matched[0][6],
			ProtocolMethod: matched[0][7],
			ServerResponse: matched[0][8],
			SendBytes:      matched[0][9],
			UserAgent:      matched[0][11],
			// browser        :,
			// system         :,
		})
	}
	return rows
}

// regularExpression runs the expression on the given data and return back the parse data
func regularExpression(MatchedCase string, resp string) [][]string {
	re := regexp.MustCompile(MatchedCase)
	if re.MatchString(resp) {
		return re.FindAllStringSubmatch(resp, -1)
	}
	return nil
}

// UploadLogs inserts the logs into database
// func UploadLogs(path string) []Logs {
// 	file, err := OpenFile(path)
// 	err = ErrorHandling(err, "error to open file", WARNING)
// 	if err != nil {
// 		log.Println(err)
// 		return []Logs{}
// 	}
// 	return ReadFile(file)
// }
