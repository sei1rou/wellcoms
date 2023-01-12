package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	//ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)

	//入力ファイル準備
	infile, err := os.Open(flag.Arg(0))
	failOnError(err)
	defer infile.Close()

	//書き込みファイル準備
	outfile, err := os.Create("./ウェルネス・コミュニケーションズ健診データ" + time.Now().Format("20060102") + ".csv")
	failOnError(err)
	defer outfile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'
	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	// writer.Comma = '\t'
	writer.Comma = ','
	writer.UseCRLF = true

	log.Print("Start\r\n")

	DrPos := 0
	DrName := "寺門　節雄"

	//　タイトル行を読みだす
	record, err := reader.Read() // １行読み出す
	for i, v := range record {
		if v == "医師名" {
			DrPos = i
			break
		}
	}
	writer.Write(record) //１行書き出す

	if DrPos == 0 {
		log.Print("「医師名」がありません。抽出パターんを確認してください。\r\n")
	}

	for {
		record, err := reader.Read() // １行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		if record[DrPos] == "" {
			record[DrPos] = DrName
		}

		//１行書き出す
		writer.Write(record)
	}
	writer.Flush()
	log.Print("Finesh !\r\n")

}
