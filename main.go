package main

import (
	"fmt"
	proc "goprod-8-dz-7/proc"
	util "goprod-8-dz-7/utils"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	numJobProcessors = 3
	numURLs          = 100
)

func main() {
	jobs := make(chan proc.Job, numURLs)
	results := make(chan proc.Result, numURLs)
	var wg sync.WaitGroup

	for w := 0; w < numJobProcessors; w++ {
		wg.Add(1)
		go proc.JobVerifyProcesor(jobs, results, &wg, w+1)
	}

	//  принимаем URL для обработки, запаковываем их в задания и кладем в очередь на обработку
	for i := 0; i < numURLs; i++ {
		var ip = net.IPv4(byte(10+rand.Intn(246)), byte(10+rand.Intn(246)), byte(10+rand.Intn(246)), byte(10+rand.Intn(246)))

		var url = fmt.Sprintf("http://%s/%s/", ip.String(), util.GenerateRandomString(3))
		job := proc.Job{ID: i + 1, URL: url, Timestamp: time.Now().Format("2006-01-02 15:04:05")}
		jobs <- job

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Случайная задержка перед добавлением следующего задания
		//fmt.Println("Added job:", job.URL)
	}
	close(jobs) // Закрываем канал с заданиями, чтобы горутины  воркеры знали, когда остановиться.

	// Запускаем отдельную горутину для ожидания завершения воркеров и закрытия канала results
	go func() {
		wg.Wait()      // Ждем завершения всех горутин-обработчиков очереди results
		close(results) // Закрываем канал с результатами
	}()

	// Выводим результаты
	// Финальный отчёт: список всех обработанных URL с их временем выполнения и общую статистику (например, среднее время, количество успешных операций).
	var totalDuration time.Duration = 0
	var totalCnt int = 0
	var successfulJobs int = 0

	for result := range results {
		fmt.Printf("● Job ID: %d, URL: %s, Status: %s, dur: %v \n", result.Job.ID, result.Job.URL, result.Status, result.Duration)
		totalDuration += result.Duration
		totalCnt++
		if result.Status == "Ok" {
			successfulJobs++
		}
	}
	fmt.Printf("Общее время: %v sec. %v url\n", totalDuration, totalCnt)
	fmt.Printf("из них количество успешных ответов: %d\n", successfulJobs)
}
