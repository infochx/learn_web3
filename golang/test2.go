package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	// ========指针=================================================================================================================================
	// 1. 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
	// 指针的使用、值传递与引用传递的区别。
	num := 0
	goPointPlus(&num)
	fmt.Println(num)

	// 2. 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
	// 指针运算、切片操作。
	arr := []int{1, 2, 3}
	slicesPlus(&arr)
	fmt.Println(arr)

	// =======Goroutine==================================================================================================================================
	// 1. 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	// go 关键字的使用、协程的并发执行。
	Goroutine()

	// 2. 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	// 协程原理、并发任务调度。
	schddule()

	// =======面向对象==================================================================================================================================
	// 1. 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	// 接口的定义与实现、面向对象编程风格。
	rectangle := Rectangle{}
	rectangle.Area()
	rectangle.Perimeter()

	circle := Circle{}
	circle.Area()
	circle.Perimeter()

	// 2. 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
	// 组合的使用、方法接收者。
	emp := Employee{person: Person{Name: "zhangsan", Age: 20}, EmployeeID: "001"}
	emp.PrintInfo()

	// =======Channel==================================================================================================================================
	// 1. 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	// 通道的基本使用、协程间通信。
	// ChannelCommunicate()

	// 2. 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	// 通道的缓冲机制。
	// ChannelMessage()

	// ==========锁机制===============================================================================================================================
	// 1. 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	// sync.Mutex 的使用、并发数据安全。
	CountWithMutex()
	// 2. 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	// 原子操作、并发数据安全。
	CountWithAtomic()
}

func CountWithAtomic() {
	var count int64 = 0
	wg := sync.WaitGroup{}
	wg.Add(10)
	start := time.Now()
	for range 10 {
		go func() {
			defer wg.Done()
			for range 1000 {
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("CountWithAtomic is finish, the count is %d and time used is %v \n", count, time.Since(start).Microseconds())
}

func CountWithMutex() {
	count := 0
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(10)
	start := time.Now()
	for range 10 {
		go func() {
			defer wg.Done()
			for range 1000 {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("CountWithMutex is finish, the count is %d and time used is %v \n", count, time.Since(start).Microseconds())
}

func ChannelMessage() {
	ch := make(chan int, 10)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(cha chan<- int) {
		defer wg.Done()
		for i := 1; i < 101; i++ {
			cha <- i
			fmt.Println("向缓冲通道发送了：", i)
		}
		close(cha)
	}(ch)

	go func(cha <-chan int) {
		defer wg.Done()
		for val := range cha {
			fmt.Println("从缓冲通道接收到：", val)
		}
	}(ch)
	wg.Wait()
}

func ChannelCommunicate() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(cha chan<- int) {
		defer wg.Done()
		for i := 1; i < 10; i++ {
			cha <- i
			fmt.Println("发送数据：", i)
		}
		close(cha)
	}(ch)

	// go func(cha <-chan int) {
	// 	defer wg.Done()
	// 	for val := range cha {
	// 		fmt.Println("接收到数据：", val)
	// 	}
	// }(ch)

	go func() {
		defer wg.Done()
		for {
			val, ok := <-ch
			if !ok {
				fmt.Println("未接收到数据！")
				break
			}
			fmt.Println("接收到数据：", val)
		}
	}()
	wg.Wait()
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	person     Person
	EmployeeID string
}

func (emp *Employee) PrintInfo() Employee {
	fmt.Println("Employee' name is", emp.person.Name, ", age is", emp.person.Age, ", employeeId is", emp.EmployeeID)
	return *emp
}

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
}

func (shape *Rectangle) Area() {
	fmt.Println("this is Rectangle's method of Area")
}
func (shape *Rectangle) Perimeter() {
	fmt.Println("this is Rectangle's method of Perimeter")
}

type Circle struct {
}

func (shape *Circle) Area() {
	fmt.Println("this is Circle's method of Area")
}
func (shape *Circle) Perimeter() {
	fmt.Println("this is Circle's method of Perimeter")
}

// 2. 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
func schddule() {
	task1 := func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("任务1完成")
	}

	task2 := func() {
		time.Sleep(1000 * time.Microsecond)
		fmt.Println("任务2完成")
	}

	task3 := func() {
		time.Sleep(10000000 * time.Nanosecond)
		fmt.Println("任务3完成")
	}

	task4 := func() {
		time.Sleep(1 * time.Second)
		fmt.Println("任务4完成")
	}

	tasks := []func(){task1, task2, task3, task4}

	wg := sync.WaitGroup{}
	startTime := time.Now()
	for i, task := range tasks {
		wg.Add(1)
		go func(id int, t func()) {
			defer wg.Done()
			taskStart := time.Now()
			fmt.Printf("任务%d开始执行\n", id+1)
			t()
			fmt.Printf("任务%d耗时:%v\n", id+1, time.Since(taskStart))
		}(i, task)
	}
	wg.Wait()
	fmt.Printf("任务总耗时:%v\n", time.Since(startTime))
}

// 1. 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func Goroutine() {
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		for i := 1; i < 10; i++ {
			if i&1 == 1 {
				fmt.Printf("奇数：%d\t", i)
			}
		}
		fmt.Println()
	}()

	wg.Add(1)
	go func() {
		wg.Done()
		for i := 2; i <= 10; i++ {
			if i&1 == 0 {
				fmt.Printf("偶数：%d\t", i)
			}
		}
		fmt.Println("")
	}()
	wg.Wait()

}

// 指针：2.实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func slicesPlus(arrPo *[]int) {
	for index := range *arrPo {
		(*arrPo)[index] *= 2
	}
}

// 指针：1、编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10
func goPointPlus(num *int) {
	*num += 10
}
