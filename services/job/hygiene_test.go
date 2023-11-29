package job

import (
	"bufio"
	"encoding/csv"
	"os"
	"testing"
)

func TestProcess(t *testing.T) {
	input := "/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/full_20231107-030703_Imp_n_click_data.csv.source"
	output := "/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/full_20231107-030703_Imp_n_click_data.csv"
	if err := processCSVFile(input, output, "gs://lr-select-vm-us-qa-temp/721211/in/inp-clid/full_20231107-030703_Imp_n_click_data.csv"); err != nil {
		t.Fatal(err)
	}
}

// Xi2970ep4re93fXyoORdLeeeit34Ob0iLetPmzLv4e17jLoA4KgUe6U7Xt1CMrfotPL1Cy|Fashmob Bandana Prints Maxi Dress|csa_refash_liv

func TestReding(t *testing.T) {

	// input := "/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/full_20231107-030703_Imp_n_click_data.csv.source"
	// output := "/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/full_20231107-030703_Imp_n_click_data.csv"
	// if err := processCSVFile(input, output); err != nil {
	// 	logs.Error("Hygiene: remove quotes failed.", err)
	// }

	file, _ := os.Open("/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/full_20231107-030703_Imp_n_click_data.csv.source1")
	outputPath := "/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/full_20231107-030703_Imp_n_click_data.csv.source"
	// 打开输出文件
	outputFile, _ := os.Create(outputPath)

	defer outputFile.Close()
	defer file.Close()
	reader := bufio.NewReader(file)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// 使用 Scanner 逐行扫描文件
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	var records [][]string

	for i := 0; i < 10000; i++ {
		records, _ = readBatch(scanner)
		for _, record := range records {
			writer.Write(record)
		}
	}

	// for i := 0; i < 1000 && scanner.Scan(); i++ {
	// 	record := parseCSVLine(scanner.Text())
	// 	records = append(records, record)
	// }
	// for i := 0; i < 1000 && scanner.Scan(); i++ {
	// 	record := parseCSVLine(scanner.Text())
	// 	records = append(records, record)
	// }
	// records, _ := readBatch(reader)
	// records, _ = readBatch(scanner)
	// for _, record := range records {
	// 	// 写入处理后的记录到输出文件
	// 	writer.Write(record)
	// }
	// records, _ = readBatch(scanner)
	// for _, record := range records {
	// 	writer.Write(record)
	// }
}
