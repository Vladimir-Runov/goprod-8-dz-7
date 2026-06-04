package main

// cd путь_кпапке
// go work init
// go mod init goprod-8-dz-7
//
// https://pkg.go.dev/golang.org/x/tools/internal/typesinternal#BrokenImport

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

type Job struct {
	ID        int
	URL       string // IP адрес клиента
	timestamp string // время в формате "2024-01-15 10:30:00"
}

// Определите структуру Result для результата (например, с полями Job, Status, Duration).
type Result struct {
	Job      Job
	Status   string
	Duration time.Duration
}

const (
	numJobProcessors = 3
	numURLs          = 100
)

var charset0 = "abcdefghijklmnopqrstuvwxyz"

// generateRandomString создает строку случайной длины из букв и символов.
func generateRandomString(length int) string {
	//	const charset1 = "abcdefghijklmnopqrstuvwxyz"
	//var charset = shuffleString()
	var sb strings.Builder // Исправлено здесь
	for i := 0; i < length; i++ {
		sb.WriteByte(charset0[rand.Intn(len(charset0))])
	}
	return sb.String()
}

// asyncJobProcesor - функция-обработчик, которая читает задания из канала jobs, обрабатывает их и отправляет результаты в канал results.
//
//	Функция должна принимать на вход канал для чтения заданий (<-chan Job), канал для записи результатов (chan<- Result) и *sync.WaitGroup.
func asyncJobProcesor(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		var sleepDuration = time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(sleepDuration) // Симуляция времени обработки

		results <- Result{
			Job:      job,
			Status:   "Processed",
			Duration: sleepDuration,
		}
	}
}
func main() {
	jobs := make(chan Job, numURLs)
	results := make(chan Result, numURLs)
	var wg sync.WaitGroup

	for w := 0; w < numJobProcessors; w++ {
		wg.Add(1)
		go asyncJobProcesor(jobs, results, &wg)
	}

	//  принимаем URL для обработки, запаковываем их в задания и кладем в очередь на обработку
	for i := 0; i < numURLs; i++ {
		var ip = net.IPv4(byte(10+rand.Intn(246)), byte(10+rand.Intn(246)), byte(10+rand.Intn(246)), byte(10+rand.Intn(246)))

		var url = fmt.Sprintf("http://%s/%s/", ip.String(), generateRandomString(3))
		job := Job{ID: i + 1, URL: url, timestamp: time.Now().Format("2006-01-02 15:04:05")}
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
	var successfulJobs int = 0

	for result := range results {
		fmt.Printf("● Job ID: %d, URL: %s, Status: %s, dur: %v \n", result.Job.ID, result.Job.URL, result.Status, result.Duration)
		totalDuration += result.Duration
		if result.Status == "Processed" {
			successfulJobs++
		}
	}
	fmt.Printf("Общее время: %v sec.\n", totalDuration)
	fmt.Printf("Количество успешных задач: %d\n", successfulJobs)
}
