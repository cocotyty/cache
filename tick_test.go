package cache

import (
	"testing"
	"time"
	"fmt"
)

func TestTicker(t *testing.T) {
	ticker:=time.NewTicker(10*time.Second)
	go func() {
		<-time.After(13*time.Second)
		ticker.Stop()
		fmt.Println("closed")
	}()
	<-ticker.C
	fmt.Println("tick")
	<-ticker.C
	fmt.Println("tick2")

}