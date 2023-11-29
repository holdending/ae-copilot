package job

import (
	"bufio"
	"encoding/csv"
	"os"
	"path"
	"strings"

	"github.com/LiveRamp/ae-copilot/config"
	"github.com/LiveRamp/ae-copilot/models"
	"github.com/LiveRamp/ae-copilot/pkg/libs/storage"
	constant "github.com/LiveRamp/ae-copilot/utils"
	"github.com/astaxie/beego/logs"
)

const (
	batchSize = 1000 // 每批次的记录数
)

type Hygiene struct {
}

func NewHygiene() *Hygiene {
	return &Hygiene{}
}

func (h *Hygiene) Running(task *models.RejectedFileRemediationTask) error {
	logs.Info("Hygiene: start to running.")

	if err := h.doing(task); err != nil {
		return err
	}
	logs.Info("Hygiene: finished task", task.TaskName)
	return nil
}
func (h *Hygiene) doing(task *models.RejectedFileRemediationTask) error {
	fs := storage.NewStorageClient(task.RejectedPrefix, config.Agent.GCSCredentials)
	fileName := path.Base(task.RejectedPrefix)

	sourceFile := constant.TEMP_DIR + fileName + constant.DOWNLOAD_SUFFIX
	temFile := constant.TEMP_DIR + fileName
	logs.Info("Hygiene: start to download.", sourceFile)
	if err := fs.Download(task.RejectedPrefix, sourceFile); err != nil {
		logs.Error("Hygiene: download file failed.", err)
		return err
	}
	defer os.Remove(sourceFile)
	err := processCSVFile(sourceFile, temFile, task.InPrefix)
	if err != nil {
		logs.Error("Hygiene: remove quotes failed.", err)
		return nil
	}
	return nil
}

func processCSVFile(inputPath, outputPath, inPrefix string) error {
	logs.Info("Hygiene: start to process csv file.")
	// 打开原始文件
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(inputPath)

	// 打开输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	defer os.Remove(outputPath)

	// 创建 CSV 读取器和写入器
	// reader := csv.NewReader(file)
	reader := bufio.NewReader(file)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	// 逐行读取和处理
	for {
		records, err := readBatch(scanner)
		if err != nil {
			logs.Error("Hygiene: read batch failed.", err)
			break // 文件读取出错，跳出循环
		}
		if len(records) == 0 {
			break // 文件读取完毕,跳出循环
		}

		// 处理每行记录的逻辑，例如删除双引号
		for _, record := range records {
			for i, field := range record {
				record[i] = removeQuotesFromString(field)
			}
			// 写入处理后的记录到输出文件
			err = writer.Write(record)
			if err != nil {
				logs.Error("Hygiene: write file failed.", err)
				return err
			}
		}
	}
	logs.Info("Hygiene: start to upload csv file.")
	fs := storage.NewStorageClient(inPrefix, config.Agent.GCSCredentials)
	return fs.Upload(outputPath, inPrefix)
}

func readBatch(scanner *bufio.Scanner) ([][]string, error) {
	var records [][]string
	// 使用 Scanner 逐行扫描文件
	for i := 0; i < batchSize && scanner.Scan(); i++ {
		record := parseCSVLine(scanner.Text())
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func parseCSVLine(line string) []string {
	// 使用 CSV 包的 Reader 读取一行记录
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(line)))
	record, _ := reader.Read()
	return record
}

func removeQuotesFromString(input string) string {
	// 移除双引号
	return strings.Replace(input, "\"", "", -1)
}
