package main

import "context"

type Job struct {
	Id     int
	Number int
	Result int
}

func worker(fn func(int) int, jobs <-chan Job, results chan<- Job) {
	for job := range jobs {
		job.Result = fn(job.Number)
		results <- job
	}
}

func ParallelMapCtx(ctx context.Context, numbers []int, fn func(int) int, maxWorkers int) ([]int, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	numberJobs := len(numbers)

	jobs := make(chan Job, numberJobs)
	results := make(chan Job, numberJobs)

	defer close(jobs)

	for range maxWorkers {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			go worker(fn, jobs, results)
		}
	}

	for id, number := range numbers {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			jobs <- Job{Id: id, Number: number}
		}
	}

	numbersResult := make([]int, numberJobs)
	for range numberJobs {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case job := <-results:
			numbersResult[job.Id] = job.Result
		}
	}

	return numbersResult, nil
}
