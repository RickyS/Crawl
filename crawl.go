//  coding: utf-8
// crawl.go by Ricky Seltzer rickyseltzer@gmail.com.  Version 2.0 on 2013-07-10

// Program to crawl the web using package creep.
package main

import (
	"creep"
	"fmt"
	"runtime"
	"time"
)

var urlCount int = 0
var urlLength int = 0            // Sum of all lengths of url strings.
var sumElapsed time.Duration = 0 // Total of elapsed time of Get() calls, even though they overlap.

var statusCodeCounts [601]int // store counts of status code appearances.  Mostly 200s, some 404s...

// Just a looping goroutine that dumps diagnostic output.
/*func monitor() {
	var sleeper int = 2 // number of seconds to nap.  To start.
	var counter int = 0
	var avgLen float64 = 0.0
	var avgDur time.Duration = 0

	for {
		synched.Lock() // Write lock shared data.
		qC := synched.queueCnt
		rJ := synched.rejectCnt
		uC := urlCount
		uL := urlLength
		synched.Unlock() // Write unlock shared data.

		if uC > 0 {
			avgLen = float64(uL) / float64(uC)
			avgDur = sumElapsed / time.Duration(uC)
		}

		strange := (qC - reqChanCapacity) == (rJ + goingCount) // true when it used to hang.  Really odd...
		straTF := boolTF(strange)
		rr := log.Sprintf("%6d:%4d.", len(reqChan), len(respChan))
		log.Printf("Gos: %3d, req:resp %s, Urls %4d, enQ %4d, avgLen %2.2f, avgDur %12v. stra %s, ml %4d\n",
			runtime.NumGoroutine(), rr, urlCount, qC, avgLen, avgDur, straTF, mapLength())
		log.Println("go Status: ", strings.Join(routineStatus, ""))
		counter++
		if 10 > counter {
			sleeper += 2
		} else if 20 > counter {
			sleeper += 5
		}
		time.Sleep(time.Duration(sleeper) * time.Second)
	}
}*/

func main() {
	jobData := creep.LoadJobData("iana.json") // Load job request into struct jobData.
	//go monitor()

	startTime := time.Now()
	for i := 0; i < len(jobData.Tests); i++ { // Once for each test in the jobData
		eachGroup := jobData.Tests[i]
		fmt.Println("")
		testnameDisplay := "'" + eachGroup.Testname + "'"
		fmt.Printf("Job %12s has: Maxurls %3d, Gomaxprocs %2d, MaxGoRoutines %3d, ExpectFail %v, JustOneDomain %s, %2d urls:",
			testnameDisplay, eachGroup.Maxurls, eachGroup.Gomaxprocs, eachGroup.MaxGoRoutines,
			eachGroup.ExpectFail, boolTF(eachGroup.JustOneDomain), len(eachGroup.Urls))

		runtime.GOMAXPROCS(eachGroup.Gomaxprocs)

		for urlNum := 0; urlNum < len(eachGroup.Urls); urlNum++ { // Once for each url in the current test.
			fmt.Printf("TestUrl # %3d: %s\n", urlNum, eachGroup.Urls[urlNum])
		}

		doJob(&eachGroup)
	}
	fmt.Printf("Test ending after %v simulating %v\n", time.Since(startTime), sumElapsed)
}

func doJob(pEachGroup *creep.JobData) {

	urls := pEachGroup.Urls
	expectFail := pEachGroup.ExpectFail
	maxurls := pEachGroup.Maxurls
	maxGoRoutines := pEachGroup.MaxGoRoutines
	testname := pEachGroup.Testname
	testnameDisplay := "'" + testname + "'"

	respChan := creep.CreepWebSites(urls, maxurls, maxGoRoutines, pEachGroup.JustOneDomain) // Call the software under test (SUT)

OnceForEachResponse:
	for {
		result, notDoneYet := <-respChan
		if !notDoneYet {
			fmt.Printf("Job Closed %s\n", testnameDisplay)
			//Channel has been closed by waitGroup, we should be all done by now.
			// fmt.Printf("Job Closed %12s: %4d urls Fetched, %4d dupes. Elapsed: %v, len (reqQ) %3d, resps: %3d\n\n",
			// 	testnameDisplay, synched.urlsFetched, synched.dupsStopped, sumElapsed, len(reqChan), urlCount)
			// fmt.Println("go Status: ", strings.Join(routineStatus, ""))
			ShowSummary()
			return
		}

		if (nil != result) && ("DONE" == result.Url) {
			fmt.Printf("Job Done %s\n", testnameDisplay)
			//Channel has been closed, we should be all done by now.
			// testnameDisplay := "'" + testname + "'"
			// fmt.Printf("Job Done %12s: %4d urls Fetched, %4d dupes. Elapsed: %v, len (reqQ) %3d, resps: %3d\n\n",
			// 	testnameDisplay, synched.urlsFetched, synched.dupsStopped, sumElapsed, len(reqChan), urlCount)
			// fmt.Println("go Status: ", strings.Join(routineStatus, ""))
			ShowSummary()
			return
		}

		if nil == result {
			fmt.Errorf("UnExpected NIL result\n")
			return
		}

		sumElapsed += result.ElapsedTime
		urlCount++
		urlLength += len(result.Url)

		if (nil != result) && (nil != result.HttpResponse) {
			sc := result.HttpResponse.StatusCode
			if (0 <= sc) && ((len(statusCodeCounts) - 1) > sc) {
				statusCodeCounts[sc]++
			} else {
				statusCodeCounts[len(statusCodeCounts)-1]++ // Count Invalid status codes.
			}
		}
		if expectFail {
			if nil == result.Err { // no fail.  Only the package tester uses this facility.
				fmt.Errorf("\n ==>> ERR (did not expect success): for %s\n", result.Url)
				// should we print stuff out if we got a good result on a fake url??
			}
			continue OnceForEachResponse
		} else {
			var statuscode int = -1
			if (nil != result) && (nil != result.HttpResponse) {
				statuscode = result.HttpResponse.StatusCode
			}

			if nil != result.Err { // did fail on a 'good' url, not expected to fail.
				// Errors out in the web can produce an 'error', not really a program error.
				fmt.Printf("\n ==>> ERR (did not expect error): [%3d] %v\n", statuscode, result.Err)
				continue OnceForEachResponse
			}
			if statuscode == -1 {
				fmt.Printf("Test got result %d on '%s'\n", statuscode, result.Url)
			}
		}
	} // end for
}

// Display a bool as a single letter, T or F.
func boolTF(george bool) string {
	if george {
		return "T"
	} else {
		return "F"
	}
}

// Show counts of status codes.
func ShowSummary() {
	var total = 0
	for sc, sCount := range statusCodeCounts {
		if 0 < sCount {
			total += sCount
			fmt.Printf("StatusCode  %3d:  %6d\n", sc, sCount)
		}
	}
	fmt.Printf("StatusCode Total: %6d\n", total)
}
