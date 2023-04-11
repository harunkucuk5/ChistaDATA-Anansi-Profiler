package postgres

import (
	"errors"
	"fmt"
	"github.com/ChistaDATA/ChistaDATA-Profiler-for-ClickHouse.git/pkg/stucts"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var prefixMap map[string]string = map[string]string{
	"%m": `(\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2}\.\d+(?:[ \+\-][A-Z\+\-\d]{3,6})?)`, // timestamp with milliseconds
	"%p": `(\d+)`,
}

//'%a' => [('t_appname',    "(.*?)"  )],
//'%u' => [('t_dbuser',       '([0-9a-zA-Z\_\[\]\-\.]*)')],					 # user name
//'%d' => [('t_dbname',       '([0-9a-zA-Z\_\[\]\-\.]*)')],					 # database name
//'%r' => [('t_hostport',     '([a-zA-Z0-9\-\.]+|\[local\]|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|[0-9a-fA-F:]+)?[\(\d\)]*')],     # remote host and port
//'%h' => [('t_client',       '([a-zA-Z0-9\-\.]+|\[local\]|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|[0-9a-fA-F:]+)?')],	      # remote host
//'%p' => [('t_pid',	  '(\d+)')],		# process ID
//'%Q' => [('t_queryid',	  '([\-]*\d+)')],	# Query ID
//'%n' => [('t_epoch',    '(\d{10}\.\d{3})')],    # timestamp as Unix epoch
//'%t' => [('t_timestamp',    '(\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2})[Z]*(?:[ \+\-][A-Z\+\-\d]{3,6})?')],      # timestamp without milliseconds
//'%m' => [('t_mtimestamp',   '(\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2})\.\d+(?:[ \+\-][A-Z\+\-\d]{3,6})?')], # timestamp with milliseconds
//'%l' => [('t_session_line', '(\d+)')],							# session line number
//'%s' => [('t_session_timestamp', '(\d{4}-\d{2}-\d{2}[\sT]\d{2}):\d{2}:\d{2}(?:[ \+\-][A-Z\+\-\d]{3,6})?')],    # session start timestamp
//'%c' => [('t_session_id',	'([0-9a-f\.]*)')],					       # session ID
//'%v' => [('t_virtual_xid',       '([0-9a-f\.\/]*)')],					     # virtual transaction ID
//'%x' => [('t_xid',	       '([0-9a-f\.\/]*)')],					     # transaction ID
//'%i' => [('t_command',	   '([0-9a-zA-Z\.\-\_\s]*)')],					# command tag
//'%e' => [('t_sqlstate',	  '([0-9a-zA-Z]+)')],					      # SQL state
//'%b' => [('t_backend_type',	  '(.*?)')],					      # backend type

var logRegex *regexp.Regexp
var logPartsLen int
var partSymbolMap map[int]string

func SetParseLogV1Params(PostgresLogPrefix string) error {
	// Default value of PostgresLogPrefix
	if PostgresLogPrefix == "" {
		PostgresLogPrefix = "%m [%p] "
	}
	logRegex = postgresLogPrefixToLogRegex(PostgresLogPrefix)

	// Find indices of symbols in prefixMap and populate symbolPosMap
	symbolPosMap := make(map[int]string)
	logPartsLen = 0
	for symbol := range prefixMap {
		indices := findAllIndices(PostgresLogPrefix, symbol)
		for _, v := range indices {
			symbolPosMap[v] = symbol
		}
		logPartsLen += len(indices)
	}

	// Sort symbol positions in ascending order
	keys := make([]int, 0, len(symbolPosMap))
	for k := range symbolPosMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	partSymbolMap = map[int]string{}
	i := 0
	for ; i < len(keys); i++ {
		partSymbolMap[i] = symbolPosMap[keys[i]]
	}
	partSymbolMap[i] = "message"
	logPartsLen += 1

	return nil
}

func ParseLogV1(logLine string) (stucts.ExtractedLog, error) {
	var clickHouseLog stucts.ExtractedLog
	if parts := logRegex.FindStringSubmatch(logLine); parts != nil && len(parts) == logPartsLen+1 {
		for i := 1; i < len(parts); i++ {
			switch partSymbolMap[i-1] {
			case "%m":
				parseTimestampWithMilliseconds(parts[i], &clickHouseLog)
				break
			case "%p":
				parseProcessID(parts[i], &clickHouseLog)
				break
			case "message":
				clickHouseLog.Message = parts[i]
				break
			default:
				return clickHouseLog, errors.New(fmt.Sprintf("Parser not defined for the symbol : %s", partSymbolMap[i]))
			}
		}
		return clickHouseLog, nil
	}
	clickHouseLog.Message = logLine
	return clickHouseLog, nil
}

func parseProcessID(part string, s *stucts.ExtractedLog) error {
	pid, err := strconv.ParseInt(part, 0, 0)
	if err != nil {
		return err
	}
	s.ProcessID = pid
	return nil
}

func parseTimestampWithMilliseconds(part string, s *stucts.ExtractedLog) error {
	var loc *time.Location
	var err error
	if strings.Contains(part, "UTC") {
		loc, err = time.LoadLocation("UTC")
		if err != nil {
			return err
		}
	}
	logTime, err := time.ParseInLocation("2006-01-02 15:04:05.000", part[:23], loc)
	if err != nil {
		return err
	}
	s.Timestamp = logTime
	return nil
}

func postgresLogPrefixToLogRegex(postgresLogPrefix string) *regexp.Regexp {
	regexString := postgresLogPrefix
	regexString = strings.Replace(regexString, `[`, `\[`, -1)
	regexString = strings.Replace(regexString, `]`, `\]`, -1)
	for k, v := range prefixMap {
		regexString = strings.Replace(regexString, k, v, -1)
	}
	regexString = `^` + regexString + `(.*)$`
	return regexp.MustCompile(regexString)
}

func findAllIndices(str string, substr string) []int {
	var indices []int

	for i := 0; i < len(str); {
		index := strings.Index(str[i:], substr)
		if index == -1 {
			break
		}

		indices = append(indices, i+index)
		i += index + len(substr)
	}
	return indices
}
