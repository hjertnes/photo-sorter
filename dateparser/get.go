package dateparser

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/hjertnes/photo-sorter/constants"
	"github.com/rotisserie/eris"
	"github.com/xiam/exif"
	"os"
	"strings"
	"time"
)

func parseExifDate(date string) (time.Time, error){
	re := regexp2.MustCompile(`^\d\d\d\d:\d\d:\d\d \d\d:\d\d:\d\d$`, 0)

	m, err := re.MatchString(date)
	if err != nil{
		return time.Time{}, eris.Wrap(err, "failed to check with regexp")
	}

	if m {
		t, err := time.Parse("2006:01:02 15:04:05", date)
		if err != nil{
			return time.Time{}, eris.Wrap(err, "failed to parse date")
		}
		return t, nil
	}
	fmt.Println(date)

	return time.Time{}, eris.Wrap(constants.ErrNotSupported, "date format not supported")
}

func readExifDate(filename string) (time.Time, error){
	data, err := exif.Read(filename)
	if err != nil{
		return time.Time{}, eris.Wrap(err, "faild to read exif")
	}

	var data1 string
	var data2 string
	var data3 string

	for i, v := range data.Tags {
		if strings.ToLower(i) == "date and time (original)"{
			data1 = v
		}
		if strings.ToLower(i) == "date and time (digitised)"{
			data2 = v
		}

		if strings.ToLower(i) == "date and time (digitized)"{
			data3 = v
		}
	}

	if len(data2) > 0{
		t, err := parseExifDate(data2)
		if err != nil{
			return time.Time{}, eris.Wrap(err, "failed to parse date")
		}

		return t, nil
	}

	if len(data3) > 0{
		t, err := parseExifDate(data3)
		if err != nil{
			return time.Time{}, eris.Wrap(err, "failed to parse date")
		}

		return t, nil
	}

	if len(data1) > 0{
		t, err := parseExifDate(data1)
		if err != nil{
			return time.Time{}, eris.Wrap(err, "failed to parse date")
		}

		return t, nil
	}

	return time.Time{}, eris.Wrap(constants.ErrNotFound, "no date found")
}
func readStatDate(filename string) (time.Time, error){
	s, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, eris.Wrap(err, "failed to stat file")
	}

	return s.ModTime(), err
}

func GetDate(filename string) time.Time{
	date, err := readExifDate(filename)
	if err == nil{
		return date
	}

	date, err = readStatDate(filename)
	if err == nil{
		fmt.Println("File system date")
		return date
	}

	fmt.Println("Fallback date")
	return time.Now()
}
