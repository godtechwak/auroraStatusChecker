package rds_util

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func PrepareCheck(region, worktype string, sess *session.Session) (*rds.RDS, []string) {
	filePath := fmt.Sprintf("./db_%s_list_%s.txt", worktype, region)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	filteredLines := filterEmptyStrings(lines)
	svc := rds.New(sess)

	return svc, filteredLines
}

func filterEmptyStrings(lines []string) []string {
	var result []string
	for _, line := range lines {
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

type result struct {
	name string
	row  []string
}

func AuroraVersionParam(svc *rds.RDS, lines []string, w *tabwriter.Writer, num int64, describable Describable, logger *log.Logger) {
	initMap := make(map[string]time.Time) // 수행시간을 담기 위한 맵
	var mu sync.Mutex                     // 맵 접근을 위한 뮤텍스
	var count int                         // 반복횟수
	count = 1

	headers := describable.GetHeaders()
	for {
		printHeader(w, headers, logger)

		var wg sync.WaitGroup
		results := make(chan result, len(lines))

		for _, name := range lines {
			wg.Add(1)
			go func(name string) {
				defer wg.Done()
				currentTime := time.Now()
				timeString := currentTime.Format("15:04:05.000")

				var duration time.Duration
				mu.Lock()
				if count > 1 {
					firstTime := initMap[name]
					duration = currentTime.Sub(firstTime)
				} else {
					initMap[name] = currentTime
				}
				mu.Unlock()

				durationString := fmt.Sprintf("%.2f", duration.Seconds()) + "s"

				output, err := describable.Describe(svc, name)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				results <- result{
					name: name,
					row:  append([]string{timeString, durationString}, output...),
				}
			}(name)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		var sortedResults []result
		for result := range results {
			sortedResults = append(sortedResults, result)
		}

		sort.Slice(sortedResults, func(i, j int) bool {
			return sortedResults[i].name < sortedResults[j].name
		})

		for _, result := range sortedResults {
			printRow(w, result.row, logger)
		}

		printFooter(w, len(headers), logger)
		time.Sleep(time.Duration(num) * time.Millisecond)
		count++
	}
}

func printHeader(w *tabwriter.Writer, headers []string, logger *log.Logger) {
	printRow(w, headers, logger)
	printDivider(w, len(headers), logger)
}

func printFooter(w *tabwriter.Writer, colCount int, logger *log.Logger) {
	printDivider(w, colCount, logger)
	fmt.Print("\033[2J\033[H")
	w.Flush()
}

func printRow(w *tabwriter.Writer, columns []string, logger *log.Logger) {
	row := strings.Join(columns, "│\t ") + "│\t\n"
	fmt.Fprintf(w, row)
	logger.Println(strings.Join(columns, " | "))
}

func printDivider(w *tabwriter.Writer, colCount int, logger *log.Logger) {
	divider := strings.Repeat("──────────────────────────────────\t ", colCount) + "\n"
	fmt.Fprintf(w, divider)
	logger.Println(strings.Repeat("─", colCount*32))
}
