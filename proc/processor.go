package proc

import (
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID        int
	URL       string // IP адрес клиента
	Timestamp string // время в формате "2024-01-15 10:30:00"
}

// Определите структуру Result для результата (например, с полями Job, Status, Duration).
type Result struct {
	Job      Job
	Status   string // Ok, Failed
	Duration time.Duration
	ProcId   int
}

// asyncJobProcesor - функция-обработчик, которая читает задания из канала jobs, обрабатывает их и отправляет результаты в канал results.
//
//	Функция должна принимать на вход канал для чтения заданий (<-chan Job), канал для записи результатов (chan<- Result) и *sync.WaitGroup.
func JobVerifyProcesor(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, procId int) {
	defer wg.Done()
	for job := range jobs {
		var Status = "Failed"
		var sleepDuration = time.Duration(rand.Intn(1000)) * time.Millisecond
		if rand.Intn(2) == 0 {
			sleepDuration = time.Duration(rand.Intn(2000)+1000) * time.Millisecond // Увеличиваем время запроса для "битых" ссылок
			time.Sleep(sleepDuration)                                              // Симуляция времени ожидания ответа с адреса
			Status = "Failed"
		} else {
			time.Sleep(sleepDuration) // Симуляция времени ожидания ответа с адреса
			Status = "Ok"
		}

		results <- Result{
			Job:      job,
			Status:   Status,
			Duration: sleepDuration,
			ProcId:   procId,
		}
	}
}
