package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptopay-dev/yaga/workers"
	"go.uber.org/atomic"
)

type myDelayLock struct {
	stop bool
}

func (m *myDelayLock) Next(t time.Time) time.Time {
	if m.stop {
		fmt.Printf("[%s] instance shell be stopped\n", time.Now().Format("15:04:05"))
		// отправляя нулевое время, мы по факту выключаем воркер,
		// больше не сможем никогда его запустить,
		// во всяком случае я пока не придумал как
		return time.Time{}
	}

	// будем планировать на каждой следующей секунде,
	// то есть фактически мы будем запускать воркер каждую секунду
	return t.Add(time.Second)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("[%s] Hello, workers!\n", time.Now().Format("15:04:05"))

	// воркер будет запускаться каждые 5 секунд
	// пример планировщика на подобии тикера
	err := workers.New(workers.Options{
		Name:     "worker #1",
		Schedule: workers.Every(time.Second * 5),
		Handler: func() {
			fmt.Printf("[%s] worker #1 every 5 secs\n", time.Now().Format("15:04:05"))
		},
	})
	if err != nil {
		panic(err)
	}

	workers.Start()

	// воркер будет запускаться каждые 13 секунд
	// пример планировщика на подобии тикера, но из парсинга строки
	sched, err := workers.Parse("@every 13s")
	if err != nil {
		panic(err)
	}
	step := atomic.NewUint32(0)
	err = workers.New(workers.Options{
		Name:     "worker #2",
		Schedule: sched,
		Handler: func() {
			fmt.Printf("[%s] worker #2 every 13 secs: STEP=%d\n", time.Now().Format("15:04:05"), step.Inc())
		},
	})
	if err != nil {
		panic(err)
	}

	// воркер будет запускаться каждую минуту в 12 секунд
	// пример планировщика на основе крона UNIX, но отличие в дополнительном
	// первом поле обозначающем секунды
	sched, err = workers.Parse("12 */1 * * * *")
	if err != nil {
		panic(err)
	}
	err = workers.New(workers.Options{
		Name:     "worker #3",
		Schedule: sched,
		Handler: func() {
			fmt.Printf("[%s] worker #3 every minute at 12 secs\n", time.Now().Format("15:04:05"))
		},
	})
	if err != nil {
		panic(err)
	}

	// воркер будет запускаться согласно кастомному планировщику
	// пример планировщика с использованием интерфейса workers.Schedule
	delay := new(myDelayLock)
	err = workers.New(workers.Options{
		Name:     "worker #4",
		Schedule: delay,
		Handler: func() {
			if step.Load() > 4 && !delay.stop {
				fmt.Printf("[%s] worker #4: send command exit\n", time.Now().Format("15:04:05"))
				delay.stop = true
				// откладываем отмену контекста на 10 секунд
				time.AfterFunc(time.Second*10, cancel)
			}
		},
	})
	if err != nil {
		panic(err)
	}

	// ждем пока контекст не будет отменен
	<-ctx.Done()

	// отдаем команду стопить все воркеры
	workers.Stop()

	// ждем пока все воркеры остановятся
	workers.Wait()

	fmt.Printf("[%s] All workers is stopped\n", time.Now().Format("15:04:05"))
}
