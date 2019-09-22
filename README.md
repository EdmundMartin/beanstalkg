Beanstalkg is a dependency free Golang Beanstalkd client.

# Usage - Producer
```go
package main

import (
	"fmt"
	"time"

	"github.com/EdmundMartin/beanstalkg"
)

func main() {
	conn, _ := beanstalkg.NewConnection(`0.0.0.0`, 11300)
	delay := time.Duration(100) * time.Second
	id, err := conn.PutString(`something`, 100, delay, delay)
	fmt.Println(id)
	fmt.Println(err)
	id, err = conn.PutString(`something`, 100, delay, delay))
}
```
Beanstalkg provides both a PutString and PutBytes method to Put a job on a Beanstalkd queue.

# Usage - Consumer
```go
package main

import (
	"fmt"
	"time"

	"github.com/EdmundMartin/beanstalkg"
)

func main() {
	conn, _ := beanstalkg.NewConnection(`0.0.0.0`, 11300)
	delay := time.Duration(10) * time.Second
	body, _ := conn.ReserveWithTimeout(delay)
	res := conn.Release(body.ID, 1029, delay)
	stats, _ := conn.Stats()
	fmt.Println(string(stats))
}
```
