package main

import "sync"

func ExecutePipeline(jobs ...job) {

	chans := make([]chan any, len(jobs)+1)
	for i := range chans {
		chans[i] = make(chan any)
	}

	wg := &sync.WaitGroup{}
	for i, j := range jobs {
		wg.Add(1)
		go func(in, out chan any, j job) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(chans[i], chans[i+1], j)
	}

	wg.Wait()
}

func SingleHash(in, out chan any) {

}

func MultiHash(in, out chan any) {

}

func CombineResults(in, out chan any) {

}
