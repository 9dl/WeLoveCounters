package main

import (
	"bufio"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func makeRequest(url *url.URL, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url.String())

	if err := client.Do(req, resp); err != nil {
		log.Fatalf("Error making request to %s: %v", url.String(), err)
	}

	if resp.StatusCode() == fasthttp.StatusOK {
		fmt.Println("Request to", url.String(), "completed successfully.")
	} else {
		fmt.Printf("Request to %s failed with status code: %d\n", url.String(), resp.StatusCode())
	}
}

func main() {
	var profileCounterURLStr string
	var threads int
	var wg sync.WaitGroup
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Profile Counter URL: ")
	profileCounterURLStr, err := reader.ReadString('\n')
	ifError("Profile Counter URL Input", err)

	profileCounterURL, err := url.Parse(profileCounterURLStr)
	ifError("Invalid Profile Counter URL", err)

	fmt.Println("Enter number of threads: ")
	threadsStr, err := reader.ReadString('\n')

	threads, err = strconv.Atoi(threadsStr)
	ifError("Thread Count Input", err)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go makeRequest(profileCounterURL, &wg)
	}

	wg.Wait()
}

func ifError(where string, err error) {
	if err != nil {
		fmt.Printf("Error in %v: %v\n", where, err)
	}
}
