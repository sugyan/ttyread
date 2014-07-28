package main

import (
	"flag"
	"fmt"
	"github.com/sugyan/ttyread"
	"io"
	"os"
	"time"
)

func main() {
	help := flag.Bool("help", false, "usage")
	speed := flag.Float64("s", 1.0, "play speed")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	filename := flag.Arg(0)
	if filename == "" {
		filename = "ttyrecord"
	}
	err := play(filename, *speed)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func play(filename string, speed float64) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	var (
		first  = true
		prevTv ttyread.TimeVal
	)
	reader := ttyread.NewTtyReader(file)
	for {
		var data *ttyread.TtyData
		data, err = reader.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		// calc delay
		var diff ttyread.TimeVal
		if first {
			first = false
		} else {
			diff = data.TimeVal.Subtract(prevTv)
		}
		prevTv = data.TimeVal
		// wait
		time.Sleep(time.Microsecond * time.Duration(float64(diff.Sec*1000000+diff.Usec)/speed))
		// write
		os.Stdout.Write(*data.Buffer)
	}
	return nil
}
