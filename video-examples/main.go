package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	waitGroup()
}













func channels() {
	messages := make(chan string)

	go func() {
		time.Sleep(4 * time.Second)
		messages <- "ping"
	}()

	msg := <-messages
	fmt.Println(msg)
}








func channelBuffering() {
	// Here we `make` a channel of strings buffering up to
	// 2 values.
	messages := make(chan string, 2)

	// Because this channel is buffered, we can send these
	// values into the channel without a corresponding
	// concurrent receive.
	messages <- "buffered"
	messages <- "channel"

	// Later we can receive these two values as usual.
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}







// This is the function we'll run in a goroutine. The
// `done` channel will be used to notify another
// goroutine that this function's work is done.
func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	// Send a value to notify that we're done.
	done <- true
}

func channelSync() {
	// Start a worker goroutine, giving it the channel to
	// notify on.
	done := make(chan bool, 1)
	go worker(done)

	// Block until we receive a notification from the
	// worker on the channel.
	<-done
}






// This `ping` function only accepts a channel for sending
// values. It would be a compile-time error to try to
// receive on this channel.
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// The `pong` function accepts one channel for receives
// (`pings`) and a second for sends (`pongs`).
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func channelDirections() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}









func funcselect() {
	// For our example we'll select across two channels.
	c1 := make(chan string)
	c2 := make(chan string)

	// Each channel will receive a value after some amount
	// of time, to simulate e.g. blocking RPC operations
	// executing in concurrent goroutines.
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	// We'll use `select` to await both of these values
	// simultaneously, printing each one as it arrives.
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

















func timeouts() {
	// For our example, suppose we're executing an external
	// call that returns its result on a channel `c1`
	// after 2s.
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	// Here's the `select` implementing a timeout.
	// `res := <-c1` awaits the result and `<-Time.After`
	// awaits a value to be sent after the timeout of
	// 1s. Since `select` proceeds with the first
	// receive that's ready, we'll take the timeout case
	// if the operation takes more than the allowed 1s.

	receiver := make(chan string)

	go func() {
		select {
		case res := <-c1:
			receiver <- res
		case <-time.After(1 * time.Second):
			receiver <- "timeout 1"
		}
	}()



	case1 := <-receiver

	fmt.Println(case1)



	// If we allow a longer timeout of 3s, then the receive
	// from `c2` will succeed and we'll print the result.
	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}























func nonBlockingChannelOperations() {
	messages := make(chan string)
	signals := make(chan bool)

	// Here's a non-blocking receive. If a value is
	// available on `messages` then `select` will take
	// the `<-messages` `case` with that value. If not
	// it will immediately take the `default` case.
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	// A non-blocking send works similarly. Here `msg`
	// cannot be sent to the `messages` channel, because
	// the channel has no buffer and there is no receiver.
	// Therefore the `default` case is selected.
	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	// We can use multiple `case`s above the `default`
	// clause to implement a multi-way non-blocking
	// select. Here we attempt non-blocking receives
	// on both `messages` and `signals`.
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}























func closingChannels() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	// Here's the worker goroutine. It repeatedly receives
	// from `jobs` with `j, more := <-jobs`. In this
	// special 2-value form of receive, the `more` value
	// will be `false` if `jobs` has been `close`d and all
	// values in the channel have already been received.
	// We use this to notify on `done` when we've worked
	// all our jobs.
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	// This sends 3 jobs to the worker over the `jobs`
	// channel, then closes it.
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	// We await the worker using the
	// [synchronization](channel-synchronization) approach
	// we saw earlier.
	<-done
}



























func rangeOverChannels() {
	// We'll iterate over 2 values in the `queue` channel.
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// This `range` iterates over each element as it's
	// received from `queue`. Because we `close`d the
	// channel above, the iteration terminates after
	// receiving the 2 elements.
	for elem := range queue {
		fmt.Println(elem)
	}
}




















// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func worker2(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func workerPools() {

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker2(w, jobs, results)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 1; a <= 5; a++ {
		<-results
	}

}





















func rateLimiting() {
	// First we'll look at basic rate limiting. Suppose
	// we want to limit our handling of incoming requests.
	// We'll serve these requests off a channel of the
	// same name.
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	// This `limiter` channel will receive a value
	// every 200 milliseconds. This is the regulator in
	// our rate limiting scheme.
	limiter := time.Tick(200 * time.Millisecond)

	// By blocking on a receive from the `limiter` channel
	// before serving each request, we limit ourselves to
	// 1 request every 200 milliseconds.
	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}








	// We may want to allow short bursts of requests in
	// our rate limiting scheme while preserving the
	// overall rate limit. We can accomplish this by
	// buffering our limiter channel. This `burstyLimiter`
	// channel will allow bursts of up to 3 events.
	burstyLimiter := make(chan time.Time, 3)

	// Fill up the channel to represent allowed bursting.
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	// Every 200 milliseconds we'll try to add a new
	// value to `burstyLimiter`, up to its limit of 3.
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	// Now simulate 5 more incoming requests. The first
	// 3 of these will benefit from the burst capability
	// of `burstyLimiter`.
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}

























func atomicCounters() {
	// We'll use an unsigned integer to represent our
	// (always-positive) counter.
	var ops uint64

	// To simulate concurrent updates, we'll start 50
	// goroutines that each increment the counter about
	// once a millisecond.
	for i := 0; i < 50; i++ {
		go func() {
			for {
				// To atomically increment the counter we
				// use `AddUint64`, giving it the memory
				// address of our `ops` counter with the
				// `&` syntax.
				atomic.AddUint64(&ops, 1)

				// Wait a bit between increments.
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Wait a second to allow some ops to accumulate.
	time.Sleep(time.Second)

	// In order to safely use the counter while it's still
	// being updated by other goroutines, we extract a
	// copy of the current value into `opsFinal` via
	// `LoadUint64`. As above we need to give this
	// function the memory address `&ops` from which to
	// fetch the value.
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}

























func mutexes() {
	// For our example the `state` will be a map.
	var state = make(map[int]int)

	// This `mutex` will synchronize access to `state`.
	var mutex = &sync.Mutex{}

	// We'll keep track of how many read and write
	// operations we do.
	var readOps uint64
	var writeOps uint64

	// Here we start 100 goroutines to execute repeated
	// reads against the state, once per millisecond in
	// each goroutine.
	for r := 0; r < 100; r++ {
		go func() {
			total := 0
			for {

				// For each read we pick a key to access,
				// `Lock()` the `mutex` to ensure
				// exclusive access to the `state`, read
				// the value at the chosen key,
				// `Unlock()` the mutex, and increment
				// the `readOps` count.
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)

				// Wait a bit between reads.
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// We'll also start 10 goroutines to simulate writes,
	// using the same pattern we did for reads.
	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Unlock()
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Let the 10 goroutines work on the `state` and
	// `mutex` for a second.
	time.Sleep(time.Second)

	// Take and report final operation counts.
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)

	// With a final lock of `state`, show how it ended up.
	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}
















type readOp struct {
	key  int
	resp chan int
}
type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func statefulGoroutines() {
	// As before we'll count how many operations we perform.
	var readOps uint64
	var writeOps uint64

	// The `reads` and `writes` channels will be used by
	// other goroutines to issue read and write requests,
	// respectively.
	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	// Here is the goroutine that owns the `state`, which
	// is a map as in the previous example but now private
	// to the stateful goroutine. This goroutine repeatedly
	// selects on the `reads` and `writes` channels,
	// responding to requests as they arrive. A response
	// is executed by first performing the requested
	// operation and then sending a value on the response
	// channel `resp` to indicate success (and the desired
	// value in the case of `reads`).
	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	// This starts 100 goroutines to issue reads to the
	// state-owning goroutine via the `reads` channel.
	// Each read requires constructing a `readOp`, sending
	// it over the `reads` channel, and the receiving the
	// result over the provided `resp` channel.
	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := &readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// We start 10 writes as well, using a similar
	// approach.
	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := &writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	// Let the goroutines work for a second.
	time.Sleep(time.Second)

	// Finally, capture and report the op counts.
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}


























func waitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		count("sheep")
		wg.Done()
	}()

	wg.Wait()
}

func count(thing string) {
	for i := 0; i < 5; i++ {
		fmt.Println(thing, i)
		time.Sleep(500 * time.Millisecond)
	}
}
