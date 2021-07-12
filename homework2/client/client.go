package client

import (
	"time"
)

func main() {
	// cli1 := http.Client{
	// 	Timeout: 2 * time.Second,
	// }
	// cli1.Do(req)

	// cli := http.Client{}
	// ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	// req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	// cli.Do(req)
}
const (
	circuitStateClosed=1
	circuitStateHalfOpen=2
)

type Proc struct{
	total int
	circuitState int 
	successes bool
	failures int
	cicrcuitOpenedAt time.Time
}

func (p *Proc) process(order string) {
	p.total++
  
	switch p.circuitState {
	case circuitStateClosed:
	   if err := pay(order); isTemp(err) {
		  p.circuitState = circuitStateHalfOpen
		  return
	   }
  
	case circuitStateHalfOpen:
	   if p.total % 2 != 0 {
		  if err := pay(order); isTemp(err) {
			 p.successes = 0
			 p.failures++
			 if p.failures >= triesBeforeSwitch {
				p.circuitState = circuitStateOpen
				p.failures = 0
				p.circuitOpenedAt = time.Now()
			 }
		  } else {
			 p.failures = 0
			 p.successes++
			 if p.successes >= triesBeforeSwitch {
				p.circuitState = circuitStateClosed
			 }
		  }
		  return
	   }
  
  
	case circuitStateOpen:
	   if p.circuitOpenedAt.Add(tryCloseTimeout).Before(time.Now()) {
		  if err := pay(order); isTemp(err) {
			 p.circuitOpenedAt = time.Now()
			 return
		  }
  
		  p.successes = 1
		  p.failures = 0
		  p.circuitState = circuitStateHalfOpen
	   }
	}
  }