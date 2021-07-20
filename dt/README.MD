# Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/attapon-th/go-pkg/dt"
)

type DT struct {
	D  dt.Date      `json:"date"`
	Dt dt.Datetime  `json:"datetime"`
	Tm dt.Timestamp `json:"timestamp"`
}

func main() {
	data := DT{
		D:  dt.Date(time.Now()),
		Dt: dt.Datetime(time.Now()),
		Tm: dt.Timestamp(time.Now()),
	}
	fmt.Printf("D: %s, Dt: %s, Tm: %s\n", data.D.String(), data.Dt.String(), data.Tm.String())

	// Encode json
	j, _ := json.Marshal(data)
	fmt.Println(string(j))

	// Decode Json
	d := DT{}
	if err := json.Unmarshal(j, &d); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("D: %s, Dt: %s, Tm: %s\n", d.D.String(), d.Dt.String(), d.Tm.String())
}

```