package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
)

/*
classic pipeline demo write by perrynzhou@gmail.com
*/
const (
	batchSize = 8
)

/*
note:
        for range 在chan上有如下特性
                1.如果chan上有数据，则for 继续往下执行，如果chan没有数据则for 会阻塞
                2.如果chan被close了，则chan为nil,for range会退出循环。
*/
type PipeFeature struct {
	input1 chan int64
	input2 chan int64
	input3 chan int64
	done   chan struct{}
	stop   chan struct{}
}

func NewPipeFeature() *PipeFeature {
	return &PipeFeature{
		input1: make(chan int64, batchSize),
		input2: make(chan int64, batchSize),
		input3: make(chan int64, batchSize),
		done:   make(chan struct{}),
		stop:   make(chan struct{}),
	}
}
func (p *PipeFeature)Info() string {
        return fmt.Sprintf("%d",100)
}
func (p *PipeFeature) init() {
	log.Println("...init running...")
	defer close(p.input1)
	for {
		select {
		case <-p.done:
			log.Println("...init stop...")
			return
		default:
			time.Sleep(5 * time.Millisecond)
			p.input1 <- rand.Int63n(65535)
		}
	}
}
func (p *PipeFeature) stage1() {
	log.Println("...stage1 running...")
	defer close(p.input2)
	for v := range p.input1 { //will block util input1 close
		v = v - rand.Int63n(1024)
		p.input2 <- v
	}
	log.Println("stage1 done...")
}
func (p *PipeFeature) stage2() {
	log.Println("...stage2 running...")
	defer close(p.input3)
	for v := range p.input2 {
		v = v + 1
		p.input3 <- v
	}
	log.Println("stage2 done...")
}
func (p *PipeFeature) stage3() {
	log.Println("...stage3 running...")
	for v3 := range p.input3 { //will block
		v3 = v3 + rand.Int63n(100)
	}
	log.Println("stage3 done...")
}
func (p *PipeFeature) Run() {
	log.Println("start pipeline...")
	go p.init()   //order2- recv data from done and closed input1, return this function
	go p.stage1() //order 3-if input1 is closed,break for loop, and close input2 before return
	go p.stage2() //order 4-if input2 is closed ,break for range input2 and close input3 before return
	// order 5- if input3 is closed,stage3 return
	p.stage3()           //  will block util input3 closed after call stage2
	p.stop <- struct{}{} // order 6-send stop flag to stop chan before end Run function
}
func (p *PipeFeature) Stop() {
	p.done <- struct{}{} // order 1-let init function to stop
	//order 7 - already recv data from stop chan
	<-p.stop //wait for recv stop chan
	log.Println("stop pipeline...")
}
func main() {
	pipe := NewPipeFeature()
	defer pipe.Stop()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go pipe.Run()
	for {
		select {
		case <-sigs:
			log.Println("recieve stop signal")
			return
		}
	}
}
