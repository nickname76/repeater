# Repeater

A Go library, for easy creation of repeating function calls. Use it for things like checking updates and others periodic actions.

It is also supports aggregation of several repeaters in one manager (see [`MultiRepeater`](#multirepeater)) 

Documentation: https://pkg.go.dev/github.com/nickname76/repeater

*Please, **star** this repository, if you found this library useful.*

## Example usage

### StartRepeater

```Go
package main

import (
	"time"

	"github.com/nickname76/repeater"
)

func main() {
	stop := repeater.StartRepeater(time.Second, func() {
		println(1)
		time.Sleep(time.Second)
		println(2)
		time.Sleep(time.Second)
		println(3)
	})

	time.Sleep(time.Second * 10)

	println("STOPPIN")
	stop()
}

```

### MultiRepeater

```Go
package main

import (
	"time"

	"github.com/nickname76/repeater"
)

func main() {
	mr := repeater.NewMultiRepeater[int]()

	mr.StartRepeater(1, time.Second, func() {
		println(1)
	})

	mr.StartRepeater(2, time.Second, func() {
		println(2)
	})

	created := mr.StartRepeater(2, time.Second, func() {
		println(2)
	})
	if !created {
		println("oops, using same id with the second one repeater")
	}

	mr.StartRepeater(3, 0, func() {
		println(3)
		time.Sleep(time.Second)
		println(33)
	})

	time.Sleep(time.Second * 1)

	stopped := mr.StopRepeater(1)
	if stopped {
		println("repeater with id 1 is now stopped")
	}

	time.Sleep(time.Second * 2)

	println("Now stopping repeaters with ids 2 and 3")

	mr.StopAllRepeaters()

	println("All repeaters are now stopped")

}

```
