package usecase

import (
	"sync"
	"time"

	"github.com/dhquero/stress-test-go/internal/infra/repository"
)

type StressTestOutput struct {
	URL            string
	Requests       uint
	Concurrency    uint
	Timeout        time.Duration
	TotalTime      time.Duration
	NumberRequests uint
	StatusCode     map[int]int
}

type StressTestUseCase struct {
	url         string
	requests    uint
	concurrency uint
	timeout     uint
	mutex       sync.Mutex
}

func NewStressTestUseCase(url string, requests uint, concurrency uint, timeout uint) *StressTestUseCase {
	return &StressTestUseCase{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
		timeout:     timeout,
		mutex:       sync.Mutex{},
	}
}

func (s *StressTestUseCase) worker(jobs <-chan string, wg *sync.WaitGroup, stressTestOutput *StressTestOutput) {
	defer wg.Done()

	for url := range jobs {
		repository := repository.NewHTTPRepository(url, s.timeout)

		response, _ := repository.Get()

		s.mutex.Lock()
		stressTestOutput.NumberRequests++
		stressTestOutput.StatusCode[response.StatusCode]++
		s.mutex.Unlock()
	}
}

func (s *StressTestUseCase) Execute() (*StressTestOutput, error) {
	timeNow := time.Now()

	stressTestOutput := &StressTestOutput{
		URL:            s.url,
		Requests:       s.requests,
		Concurrency:    s.concurrency,
		Timeout:        time.Duration(s.timeout) * time.Second,
		NumberRequests: 0,
		StatusCode:     map[int]int{},
	}

	// for i := 0; i < int(s.requests); i++ {
	// 	repository := repository.NewHTTPRepository(s.url, s.timeout)

	// 	response, err := repository.Get()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	stressTestOutput.NumberRequests++
	// 	stressTestOutput.StatusCode[response.StatusCode]++
	// }

	var wg sync.WaitGroup
	jobs := make(chan string, s.concurrency)

	for workerNumber := 1; workerNumber <= int(s.concurrency); workerNumber++ {
		wg.Add(1)
		go s.worker(jobs, &wg, stressTestOutput)
	}

	for requestNumber := 0; requestNumber < int(s.requests); requestNumber++ {
		jobs <- s.url
	}
	close(jobs)

	wg.Wait()

	stressTestOutput.TotalTime = time.Since(timeNow)

	return stressTestOutput, nil
}
