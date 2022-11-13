package main

import "sync"

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan any)

	for _, j := range jobs {
		wg.Add(1)
		out := make(chan any)

		go func(in, out chan any, j job) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(in, out, j)

		in = out
	}

	wg.Wait()
}

func SingleHash(in, out chan any) {

}

func MultiHash(in, out chan any) {

}

func CombineResults(in, out chan any) {

}
