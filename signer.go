package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

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
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for d := range in {
		wg.Add(1)

		go func(in any, out chan any, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()

			data := fmt.Sprintf("%v", in)
			mu.Lock()
			md5Data := DataSignerMd5(data)
			mu.Unlock()

			crc32 := make(chan string)
			go func(s string, out chan string) {
				out <- DataSignerCrc32(s)
			}(data, crc32)

			out <- (<-crc32 + "~" + DataSignerCrc32(md5Data))
		}(d, out, wg, mu)
	}

	wg.Wait()
}

func MultiHash(in, out chan any) {
	wg := &sync.WaitGroup{}

	for d := range in {
		wg.Add(1)

		go func(in any) {
			defer wg.Done()

			iWg := &sync.WaitGroup{}
			mhs := make([]string, 6)

			for th := 0; th < 6; th++ {
				iWg.Add(1)

				go func(th int, in any, mhs []string) {
					defer iWg.Done()

					mhs[th] = DataSignerCrc32(fmt.Sprint(th) + fmt.Sprintf("%v", in))
				}(th, in, mhs)
			}

			iWg.Wait()
			out <- strings.Join(mhs, "")
		}(d)
	}

	wg.Wait()
}

func CombineResults(in, out chan any) {
	var res []string
	for d := range in {
		res = append(res, fmt.Sprintf("%v", d))
	}

	sort.Strings(res)
	out <- strings.Join(res, "_")
}
