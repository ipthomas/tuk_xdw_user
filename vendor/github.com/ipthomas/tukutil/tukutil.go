package tukutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"encoding/base64"

	"github.com/ipthomas/tukcnst"

	"github.com/google/uuid"
)

var (
	ServerName = ""
	SeedRoot   = "1.2.40.0.13.1.1.3542466645."
	IdSeed     = getIdIncrementSeed(5)
	CodeSystem = make(map[string]string)
	debugMode  = false
)

func init() {
	ServerName, _ = os.Hostname()
}
func SetDebugMode(debug bool) {
	debugMode = debug
}

// TemplateFuncMap returns a functionMap of tukutils for use in templates
func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"dtday":          Tuk_Day,
		"dtmonth":        Tuk_Month,
		"dtyear":         Tuk_Year,
		"mapval":         GetCodeSystemVal,
		"prettytime":     PrettyTime,
		"newuuid":        NewUuid,
		"newid":          Newid,
		"tuktime":        Time_Now,
		"simpledatetime": SimpleDateTime,
		"isafternow":     IsAfterNow,
		"duration":       GetDuration,
		"durationsince":  GetDurationSince,
		"hasexpired":     OHT_ShouldEscalate,
		"completedate":   OHT_FutureDate,
	}
}
func ReturnEncoded(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
func ReturnDecoded(s string) string {
	s = strings.ReplaceAll(s, " ", "+")

	str, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Printf("Error Decoding Base64 String : %s", err.Error())
		return ""
	}
	return string(str)
}
func SimpleDateTime() string {
	return Tuk_Year() + Tuk_Month() + Tuk_Day() + Tuk_Hour() + Tuk_Min() + Tuk_Sec()
}

// SplitXDWKey takes a string input (xdwkey) and returns the pathway and nhs id for the xdw
func SplitXDWKey(xdwkey string) (string, string) {
	var pwy string
	var nhs string
	if len(xdwkey) > 10 {
		pwy = xdwkey[:len(xdwkey)-10]
		nhs = strings.TrimPrefix(xdwkey, pwy)
	}
	return pwy, nhs
}

// SetCodeSystem takes a map input and sets the codesystem map with the input
func SetCodeSystem(cs map[string]string) {
	CodeSystem = cs
	log.Printf("Loaded %v code system key values", len(CodeSystem))
	if debugMode {
		Log(CodeSystem)
	}
}

// LoadCodeSystemFile loads the `codesystem.json` file from the `configs` folder
func LoadCodeSystemFile(codesystemFile string) error {
	file, err := os.Open(codesystemFile)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if err = json.NewDecoder(file).Decode(&CodeSystem); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Loaded %v code system key values", len(CodeSystem))
	if debugMode {
		Log(CodeSystem)
	}
	return nil
}

// GetCodeSystemVal takes a string input (key) and returns from the codesystem the string value corresponding to the input (key)
func GetCodeSystemVal(key string) string {
	val, ok := CodeSystem[key]
	if ok {
		return val
	}
	return key
}

// CreateLog checks if the log folder exists and creates it if not. It then checks for a subfolder for the current year i.e. 2022 and creates it if it does not exist. It then checks for a log file with a name equal to the current day and month and extension .log i.e. 0905.log. If it exists log output is appended to the existing file otherwise a new log file is created.
func CreateLog(log_Folder string) *os.File {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	mdir := log_Folder
	if _, err := os.Stat(mdir); errors.Is(err, fs.ErrNotExist) {
		if e2 := os.Mkdir(mdir, 0700); e2 != nil {
			log.Println(err.Error())
			return nil
		}
	}
	dir := mdir + "/" + Tuk_Year()
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		if e2 := os.Mkdir(dir, 0700); e2 != nil {
			log.Println(err.Error())
			return nil
		}
	}
	logFile, err := os.OpenFile(dir+"/"+Tuk_Month()+Tuk_Day()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Println("Using log file - " + logFile.Name())
	log.SetOutput(logFile)
	log.Println("-----------------------------------------------------------------------------------")
	return logFile
}
func MonitorApp() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		signalType := <-ch
		signal.Stop(ch)
		fmt.Println("")
		fmt.Println("***")
		fmt.Println("Exit command received. Exiting...")
		exitcode := 0
		switch signalType {
		case os.Interrupt:
			log.Println("FATAL: CTRL+C pressed")
		case syscall.SIGTERM:
			log.Println("FATAL: SIGTERM detected")
			exitcode = 1
		}
		os.Exit(exitcode)
	}()
}

// Log takes any struc as input and logs out the struc as a json string
func Log(i interface{}) {
	if debugMode {
		b, _ := json.MarshalIndent(i, "", "  ")
		log.Println(string(b))
	}
}

func OHT_ShouldEscalate(startdate time.Time, htDate string) bool {
	escalationdate := OHT_FutureDate(startdate, htDate)
	log.Printf("Escalation required - %v", time.Now().After(escalationdate))
	return time.Now().After(escalationdate)
}

// OHT_FutureDate takes a 'start date' and a period in the future as a string containing an OASIS Human Task api function eg. day(x) returns x days in the future from the `start date`. Valid periods are min(x),hour(x),day(x),month(x) and year(x)
func OHT_FutureDate(startdate time.Time, htDate string) time.Time {
	if strings.Contains(htDate, "(") && strings.Contains(htDate, ")") {
		periodstr := strings.Split(htDate, "(")[0]
		periodtime := GetIntFromString(strings.Split(strings.Split(htDate, "(")[1], ")")[0])
		log.Printf("Calculating date %v %s from %s", periodtime, periodstr, startdate.String())
		switch periodstr {
		case "min":
			return GetFutureDate(startdate, 0, 0, 0, 0, periodtime)
		case "hour":
			return GetFutureDate(startdate, 0, 0, 0, periodtime, 0)
		case "day":
			return GetFutureDate(startdate, 0, 0, periodtime, 0, 0)
		case "month":
			return GetFutureDate(startdate, 0, periodtime, 0, 0, 0)
		case "year":
			return GetFutureDate(startdate, periodtime, 0, 0, 0, 0)
		}
	}
	return startdate
}

// GetDurationSince takes a time as string input in RFC3339 format (yyyy-MM-ddThh:mm:ssZ) and returns the duration in days, hours and mins in a 'pretty format' eg '2 Days 0 Hrs 52 Mins' between the provided time and time.Now() as a string
func GetDurationSince(stime string) string {
	log.Println("Obtaining time Duration since - " + stime)
	st, err := time.Parse(time.RFC3339, stime)
	if err != nil {
		log.Println(err.Error())
		return "Not Available"
	}
	dur := time.Since(st)
	log.Printf("Duration - %v", dur.String())
	days := 0
	hrs := int(dur.Hours())
	min := int(dur.Minutes())

	if hrs > 24 {
		days = hrs / 24
		hrs = hrs % 24
	}
	daysstr := strconv.Itoa(days)
	hrsstr := strconv.Itoa(hrs)
	minstr := strconv.Itoa(min - (days * 24 * 60) - (hrs * 60))
	log.Println("Returning " + daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins")
	return daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins"
}

// GetDuration takes 2 times as string inputs in RFC3339 format (yyyy-MM-ddThh:mm:ssZ) and returns the duration in days, hours and mins in a 'pretty format' eg '2 Days 0 Hrs 52 Mins' between the provided times as a string
//
//	Example : GetDuration("2022-09-04T13:15:20Z", "2022-09-14T16:20:01Z") returns `10 Days 3 Hrs 4 Mins`
func GetDuration(dur time.Duration) string {
	days := 0
	hrs := int(dur.Hours())
	min := int(dur.Minutes())
	secs := int(dur.Seconds())
	if secs < 60 {
		log.Printf("Returning %v Secs", secs)
		return strconv.Itoa(secs) + " Secs"
	}
	if hrs > 24 {
		days = hrs / 24
		hrs = hrs % 24
	}
	daysstr := strconv.Itoa(days)
	hrsstr := strconv.Itoa(hrs)
	minstr := strconv.Itoa(min - (days * 24 * 60) - (hrs * 60))
	log.Println("Returning " + daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins")
	return daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins"
}

// IsAfterNow takes a time as a string input in RFC3339 format (yyyy-MM-ddThh:mm:ssZ) and returns true if the input time is after time.Now() and false if input time is before time.Now()
func IsAfterNow(inTime string) bool {
	log.Printf("Checking if %s is after the current time", inTime)
	it, err := time.Parse(time.RFC3339, inTime)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	now := time.Now().Local()
	log.Println("Time Now - " + now.Local().String())
	log.Println("Start Time - " + it.Local().String())
	log.Printf("Time %s IsAfter(time.Now()) = %v", inTime, now.Before(it))
	return now.Before(it)
}

// Pretty_Time_Now returns a pretty version of the current time for location Europe/London (strips everything after the `.` in Tuk_Time)
func Pretty_Time_Now() string {
	return PrettyTime(Time_Now())
}

// Time_Now returns the current time for location Europe/London.
func Time_Now() string {
	location, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Println(err.Error())
		return time.Now().String()
	}
	return time.Now().In(location).Format(time.RFC3339)
}

// PrettyTime fist splits the input based on sep =`.`, it takes index 0 of the split and replaces any T with a space then removes any trailing Z. It then splits the resulting string on sep = `+` returning index 0 of the split
func PrettyTime(time string) string {
	return strings.TrimSuffix(strings.Split(strings.TrimSuffix(strings.ReplaceAll(strings.Split(time, ".")[0], "T", " "), "Z"), "+")[0], " ")
}

// TUK_Hour returns the current hour as a 2 digit string
func Tuk_Hour() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Hour())
}

// TUK_Min returns the current minute as a 2 digit string
func Tuk_Min() string {
	return fmt.Sprintf("%02d", time.Now().Local().Minute())
}

// TUK_Sec returns the current second as a 2 digit string
func Tuk_Sec() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Second())
}

// TUK_MilliSec returns the current milliseconds as a 3 digit int
func Tuk_MilliSec() int {
	return GetIntFromString(Substr(GetStringFromInt(time.Now().Nanosecond()), 0, 3))
}

// TUK_Day returns the current day as a 2 digit string
func Tuk_Day() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Day())
}

// TUK_Year returns the current year as a 4 digit string
func Tuk_Year() string {
	return fmt.Sprintf("%d",
		time.Now().Local().Year())
}

// TUK_Month returns the current month as a 2 digit string
func Tuk_Month() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Month())
}

// NewUuid returns a random UUID as a string
func NewUuid() string {
	u := uuid.New()
	return u.String()
}

// GetStringFromInt takes a int input and returns a string of that value.
func GetStringFromInt(i int) string {
	return strconv.Itoa(i)
}

// GetIntFromString takes a string input with an integer value and returns an int of that value. If the input is not numeric, 0 is returned
func GetIntFromString(input string) int {
	i, _ := strconv.Atoi(input)
	return i
}

// Substr takes a string input and returns the rune (Substring) defined by the start and length in th start and length input values
func Substr(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}
	return string(asRunes[start : start+length])
}

// GetXMLNodeList takes an xml message as input and returns the xml node list as a string for the node input value provide
func GetXMLNodeList(message string, node string) string {
	if strings.Contains(message, node) {
		var nodeopen = "<" + node
		var nodeclose = "</" + node + ">"
		log.Println("Searching for XML Element: " + nodeopen + ">")
		var start = strings.Index(message, nodeopen)
		var end = strings.Index(message, nodeclose) + len(nodeclose)
		m := message[start:end]
		log.Println("Extracted XML Element Nodelist")
		return m
	}
	log.Println("Message does not contain Element : " + node)
	return ""
}

// PrettyAuthorInstitution takes a string input (XDS Author.Institution format) and returns a string with just the Institution name
func PrettyAuthorInstitution(institution string) string {
	if strings.Contains(institution, "^") {
		return strings.Split(institution, "^")[0] + ","
	}
	return institution
}

// PrettyAuthorPerson takes a string input (XDS Author.Person format) and returns a string with the person last and first names
func PrettyAuthorPerson(author string) string {
	if strings.Contains(author, "^") {
		authorsplit := strings.Split(author, "^")
		if len(authorsplit) > 2 {
			log.Println("Split Author " + authorsplit[1] + " " + authorsplit[2])
			return authorsplit[1] + " " + authorsplit[2]
		}
		if len(authorsplit) > 1 {
			return authorsplit[1]
		}
	}
	log.Println("Parsed Author " + author)
	return author
}

// GetFolderFiles takes a string input of the complete folder path and returns a fs.DirEntry
func GetFolderFiles(folder string) ([]fs.DirEntry, error) {
	var err error
	var f *os.File
	var fileInfo []fs.DirEntry
	f, err = os.Open(folder)
	if err != nil {
		log.Println(err.Error())
		return fileInfo, err
	}
	fileInfo, err = f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Println(err.Error())
	}
	return fileInfo, err
}
func getIdIncrementSeed(len int) int {
	return GetIntFromString(Substr(GetStringFromInt(time.Now().Nanosecond()), 0, len))
}
func IsBrokerExpression(exp string) bool {
	return strings.Contains(exp, "^^")
}

func GetTimeFromString(timestr string) time.Time {
	timestr = strings.Split(timestr, ".")[0]
	timestr = strings.Split(timestr, " +")[0]
	var err error
	var rsptime time.Time
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Println(err.Error())
		return rsptime
	}
	if !strings.Contains(timestr, "T") {
		rsptime, err = time.ParseInLocation("2006-01-02 15:04:05", timestr, loc)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		rsptime, err = time.ParseInLocation(time.RFC3339, timestr, loc)
		if err != nil {
			log.Println(err.Error())
			rsptime, err = time.ParseInLocation("2006-01-02T15:04:05Z", timestr, loc)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	return rsptime
}

func GetFutueDaysDate(startDate time.Time, days int) time.Time {
	return startDate.AddDate(0, 0, days)
}
func GetFutureDate(startDate time.Time, years int, months int, days int, hours int, mins int) time.Time {
	fdate := startDate.AddDate(years, months, days)
	fdate = fdate.Add(time.Hour * time.Duration(hours))
	return fdate.Add(time.Minute * time.Duration(mins))
}

// getErrorMessage returns the error message within
// the SOAP response or returns a generic error message
func GetErrorMessage(message string) string {
	if strings.Contains(message, "soap:Reason") {
		var start = strings.Index(message, "<soap:Reason>") + 13
		var end = strings.Index(message, "</soap:Reason>")
		var xmlmessage string = message[start:end]
		start = strings.Index(xmlmessage, ">") + 1
		end = strings.Index(xmlmessage, "</soap")
		return xmlmessage[start:end]
	}

	if strings.Contains(message, "faultstring") {
		var start = strings.Index(message, "<faultstring>") + 13
		var end = strings.Index(message, "</faultstring>")
		return message[start:end]
	}

	return "Soap error reason not found."
}

// containsError checks to see if the supplied
// message contains one of the two error tags
func ContainsError(message string) bool {
	if strings.Contains(message, "<soap:Fault>") {
		return true
	}
	if strings.Contains(message, "<faultstring>") {
		return true
	}
	return false
}

// getDocumentReturnList extracts the document
// return list from the SOAP response message
func GetDocumentReturnList(message string) string {
	if strings.Contains(message, "<return>") {
		var start = strings.Index(message, "<return>")
		var end = strings.Index(message, "</return>") + 9
		return message[start:end]
	}
	return message
}

func PrettyPrintDuration(duration time.Duration) string {
	rsp := ""
	secs := int(duration.Seconds())
	mins := secs / 60
	hrs := mins / 60
	days := hrs / 24
	hrs = hrs % 24
	mins = mins % 60
	if secs < 60 {
		return GetStringFromInt(secs) + " Secs"
	}
	if days == 1 {
		rsp = GetStringFromInt(days) + " Day " + GetStringFromInt(hrs) + " Hours " + GetStringFromInt(mins) + " Mins "
	}
	if days == 0 || days > 1 {
		rsp = GetStringFromInt(days) + " Days " + GetStringFromInt(hrs) + " Hours " + GetStringFromInt(mins) + " Mins "
	}
	return rsp
}
func DT_Day() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Day())
}
func DT_Hour() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Hour())
}
func DT_Min() string {
	return fmt.Sprintf("%02d", time.Now().Local().Minute())
}
func DT_Sec() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Second())
}
func DT_MilliSec() int {
	return GetMilliseconds()
}
func DT_yyyy_MM_dd() string {
	return DT_Year() + "-" + DT_MM_dd()
}
func DT_yyyy_MM_dd_hh() string {
	return DT_yyyy_MM_dd() + " " + DT_Hour()
}
func DT_yyyy_MM_dd_hh_mm() string {
	return DT_yyyy_MM_dd_hh() + ":" + DT_Min()
}
func DT_yyyy_MM_dd_hh_mm_SS() string {
	return DT_yyyy_MM_dd_hh_mm() + ":" + DT_Sec()
}
func DT_yyyy_MM_dd_hh_mm_SS_sss() string {
	return DT_yyyy_MM_dd_hh_mm_SS() + "." + GetStringFromInt(DT_MilliSec())
}
func DT_yyyyMMddhhmmSSsss() string {
	return DT_Year() + DT_Month() + DT_Hour() + DT_Min() + DT_Sec() + strconv.Itoa(DT_MilliSec())
}
func DT_SQL() string {
	return DT_Year() + "-" + DT_Month() + "-" + DT_Day() + "T" + DT_Hour() + ":" + DT_Min() + ":" + DT_Sec()
}
func DT_SQL_Future_Year() string {
	t := time.Now().Local().AddDate(1, 0, 0)
	return GetStringFromInt(t.Local().Year()) + "-" + GetStringFromInt(int(t.Local().Month())) + "-" + GetStringFromInt(t.Local().Day()) + "T" + GetStringFromInt(t.Local().Hour()) + ":" + GetStringFromInt(t.Local().Minute()) + ":" + GetStringFromInt(t.Local().Second())
}
func DT_Zulu() string {
	return DT_SQL() + "." + GetStringFromInt(DT_MilliSec()) + "Z"
}
func DT_Zulu_Future(future int64) string {
	t := time.Now().Local().Add(time.Hour*time.Duration(0) + time.Minute*time.Duration(future) + time.Second*time.Duration(0))
	return GetStringFromInt(t.Local().Year()) + "-" + GetStringFromInt(int(t.Local().Month())) + "-" + GetStringFromInt(t.Local().Day()) + "T" + GetStringFromInt(t.Local().Hour()) + ":" + GetStringFromInt(t.Local().Minute()) + ":" + GetStringFromInt(t.Local().Second()) + "." + GetStringFromInt(GetMilliseconds()) + "Z"
}
func DT_Zulu_Future_Year() string {
	t := time.Now().Local().AddDate(1, 0, 0)
	return GetStringFromInt(t.Local().Year()) + "-" + GetStringFromInt(int(t.Local().Month())) + "-" + GetStringFromInt(t.Local().Day()) + "T" + GetStringFromInt(t.Local().Hour()) + ":" + GetStringFromInt(t.Local().Minute()) + ":" + GetStringFromInt(t.Local().Second()) + "." + GetStringFromInt(GetMilliseconds()) + "Z"
}
func DT_Kitchen() string {
	return time.Now().Format(time.Kitchen)
}
func DT_Unix() string {
	return time.Now().Format(time.UnixDate)
}
func DT_ANSIC() string {
	return time.Now().Format(time.ANSIC)
}
func DT_Stamp() string {
	return time.Now().Format(time.Stamp)
}
func DT_Date() string {
	return time.Now().Format("Jan 2 2006")
}
func DT_Time() string {
	return time.Now().Format("15:04:05")
}
func DT_EPOCH() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintln(timestamp)
}
func GetMilliseconds() int {
	return GetIntFromString(Substr(GetStringFromInt(time.Now().Nanosecond()), 0, 3))
}

func Newdatetime() string {
	return DT_yyyy_MM_dd_hh_mm_SS()
}
func Newyearfuturezulu() string {
	return DT_Zulu_Future_Year()
}
func Newzulu() string {
	return DT_Zulu()
}
func New30mfutureyearzulu() string {
	return DT_Zulu_Future(30)
}
func DT_MM_dd() string {
	return DT_Month() + "-" + DT_Day()
}

func DT_Year() string {
	return fmt.Sprintf("%d",
		time.Now().Local().Year())
}
func GetIdIncrementSeed(len int) int {
	return GetIntFromString(Substr(GetStringFromInt(time.Now().Nanosecond()), 0, len))
}
func DT_Month() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Month())
}
func Loadjsonfile(jfp string) *json.Decoder {
	file, err := os.Open(jfp)
	if err != nil {
		log.Panic(jfp + " Unable to load file:-" + jfp + " - " + err.Error())
	}
	return json.NewDecoder(file)
}
func GetXdwConfigFiles(basepath string) (map[string][]byte, error) {
	var xdwFiles = make(map[string][]byte)
	var err error
	var f *os.File
	var fileInfo []fs.DirEntry
	f, err = os.Open(basepath + "xdwconfig/")
	if err != nil {
		log.Println(err.Error())
		return xdwFiles, err
	}
	fileInfo, err = f.ReadDir(-1)
	defer f.Close()
	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".json") && strings.Contains(file.Name(), "_xdwdef") {
			xdwfile, err := os.ReadFile(basepath + "xdwconfig/" + file.Name())
			if err != nil {
				log.Println(err.Error())
				return xdwFiles, err
			}
			log.Println("Loaded WF Def for Pathway : " + file.Name())
			xdwFiles[file.Name()] = xdwfile
		}
	}

	return xdwFiles, err
}
func GetHTMLWidgetFiles(basepath string) ([]string, error) {
	var htmlWidgets []string
	var err error
	var f *os.File
	var fileInfo []fs.DirEntry
	f, err = os.Open(basepath + "templates/html/")
	if err != nil {
		log.Println(err.Error())
		return htmlWidgets, err
	}
	if fileInfo, err = f.ReadDir(-1); err != nil {
		log.Println(err.Error())
		return htmlWidgets, err
	}
	defer f.Close()
	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".json") && strings.Contains(file.Name(), "_xdwdef") {
			tmplt, err := os.ReadFile(basepath + "templates/html/" + file.Name())
			if err != nil {
				log.Println(err.Error())
				return htmlWidgets, err
			}
			log.Println("Loaded html template : " + file.Name())
			htmlWidgets = append(htmlWidgets, string(tmplt))
		}
	}
	return htmlWidgets, nil
}
func Minus(n1 int, n2 int) string {
	log.Printf("Template called minus function(%v,%v) Returning %v", n1, n2, n1-n2)
	return GetStringFromInt(n1 - n2)
}

func VALUE_LIKE(val string) string {
	return "%" + val + "%"
}
func UploadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	r.ParseMultipartForm(10)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Fprintf(w, "<h3 style='color:red'>Failed to Upload File : "+err.Error()+" : </h3>")
		log.Println(err.Error())
		return
	}
	defer file.Close()
	var nhsid = r.FormValue("nhs")
	var pathway = r.FormValue("pathway")
	if nhsid == "" {
		fmt.Fprintf(w, "<h3 style='color=red'>No Patient NHS ID provided. Upload Terminated</h3>")
		return
	}
	if pathway == "" {
		fmt.Fprintf(w, "<h3 style='color=red'>No Pathway provided. Upload Terminated</h3>")
		return
	}
	fmt.Fprintf(w, "<h3>Uploading File</h3>")
	fmt.Fprintf(w, "<h3 style='color:green'>Uploading File</h3>")
	log.Printf("Uploading File: %+v", handler.Filename)
	fmt.Fprintf(w, "<h3 style='color:green'>File Size: %+v", handler.Size)
	log.Printf("File Size: %+v", handler.Size)
	log.Printf("MIME Header: %+v", handler.Header)

	fn := "uploads/" + pathway + nhsid + "_" + handler.Filename
	uploadFile, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(w, "<h3 style='color:red'>Failed to create File on Server : "+err.Error()+" : </h3>")
		log.Println(err.Error())
		return
	}
	defer uploadFile.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "<h3 style='color:red'>Failed to read file from client : "+err.Error()+" : </h3>")
		log.Println(err.Error())
		return
	}
	uploadFile.Write(fileBytes)
	fmt.Fprintf(w, "<h3 style='color:green'>Successfully Uploaded File</h3>")
	log.Printf("Successfully Uploaded File: %+v", handler.Filename)
	log.Println("Saved file to " + fn)
}
func WriteFileToTempFolder(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}
	return nil
}
func WriteResponseHeaders(fn http.HandlerFunc, secure bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Tiani_Spirit_UK")
		if r.Header.Get(tukcnst.ACCEPT) == tukcnst.APPLICATION_XML {
			w.Header().Set(tukcnst.CONTENT_TYPE, tukcnst.APPLICATION_XML)
		} else if r.Header.Get(tukcnst.ACCEPT) == tukcnst.APPLICATION_JSON {
			w.Header().Set(tukcnst.CONTENT_TYPE, tukcnst.APPLICATION_JSON)
		} else if r.Header.Get(tukcnst.ACCEPT) == tukcnst.ALL {
			w.Header().Set(tukcnst.CONTENT_TYPE, tukcnst.TEXT_HTML)
		} else {
			w.Header().Set(tukcnst.CONTENT_TYPE, tukcnst.TEXT_HTML)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "accept, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if secure {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}
		fn(w, r)
	}
}

// returns unique id in format '1.2.40.0.13.1.1.3542466645.20211021090059143.32643'
// idroot constant - 1.2.40.0.13.1.1.3542466645.
// + datetime	   - 20211021090059143.
// + 5 digit seed  - 32643
// The seed is incremented after each call to newid().
func Newid() string {
	id := SeedRoot + DT_yyyyMMddhhmmSSsss() + "." + GetStringFromInt(IdSeed)
	IdSeed = IdSeed + 1
	return id
}
func GetServiceUrl(port int, scheme, host, url string) string {
	return scheme + "://" + host + ":" + strconv.Itoa(port) + "/" + url
}
func GetGlypicon(tasktype string) string {
	switch strings.ToUpper(tasktype) {
	case "CPIS":
		return "fa fa-medkit fa-2x"
	case "TOC":
		return "fa fa-ambulance fa-2x"
	case "Transport":
		return "fa fa-ambulance fa-2x"
	default:
		return "fa fa-user-md fa-2x"
	}
}
func GetXMLNodeVal(message string, node string) string {
	if strings.Contains(message, node) {
		var nodeopen = "<" + node + ">"
		var nodeclose = "</" + node + ">"
		log.Println("Searching for value in : " + nodeopen + nodeclose)
		var start = strings.Index(message, nodeopen) + len(nodeopen)
		var end = strings.Index(message, nodeclose)
		m := message[start:end]
		log.Println("Returning value : " + m)
		return m
	}
	log.Println("Message does not contain Node : " + node)
	return ""
}
func ArrayContains(strs []string, str string) (int, bool) {
	if len(strs) > 0 {
		for s := range strs {
			if strs[s] == str {
				return s, true
			}
		}
	}
	return -1, false
}
func GetFileBytes(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	return byteValue, nil
}
func GetXmlReturnNode(message string) string {
	log.Println("Searching for <return> node in response message")
	if strings.Contains(message, "<return>") {
		var start = strings.Index(message, "<return>")
		var end = strings.Index(message, "</return>") + 9
		log.Println("Found Node <return>")
		return message[start:end]
	}
	log.Println("Node <return> Not found. Returning message")
	return message
}
func NotEmpty(params []string) bool {
	for param := range params {
		if params[param] == "" {
			return false
		}
	}
	return true
}
func SplitFhirOid(oid string) string {
	if !strings.Contains(oid, ":") {
		return oid
	}

	splitoid := strings.Split(oid, ":")
	if len(splitoid) > 2 {
		return splitoid[2]
	}
	return oid
}
func SplitExpression(exp string) string {
	if !strings.Contains(exp, "^^") {
		return exp
	}
	str := strings.Split(exp, "^^")[0]
	return str
}
