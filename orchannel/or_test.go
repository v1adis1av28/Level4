package orchannel

import (
	"sync"
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	//Ожидаем что функция отработает штатно на закрытие канала
	t.Run("closes when one of the channels closes", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		orDone := Or(ch1, ch2)

		go func() {
			time.Sleep(10 * time.Millisecond)
			close(ch1)
		}()

		select {
		case <-orDone:
		case <-time.After(50 * time.Millisecond):
			t.Fatal("Or function didn't close the channel in time")
		}
	})

	//Если в функцию не переданно каналов то должна вести корректно, и также закрывать каналы
	t.Run("closes immediately if no channels provided", func(t *testing.T) {
		orDone := Or()
		if orDone == nil {
			return
		}
		select {
		case <-orDone:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Or function didn't close with empty input")
		}
	})

	//Проверка на то как обработает уже закрытый канал
	t.Run("closes with one closed channel", func(t *testing.T) {
		ch1 := make(chan interface{})
		close(ch1)
		orDone := Or(ch1)

		select {
		case <-orDone:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Expected orDone to be closed")
		}
	})

	//Проверка что на нескольких каналах при одном из закрытых тоже сразу закроет канал
	t.Run("closes with two channels, first closes", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		close(ch1)
		orDone := Or(ch1, ch2)

		select {
		case <-orDone:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Expected orDone to be closed")
		}
	})

	//Аналогично что и тест с первым закрытым, но проверяем что и на второй закрытый он также среагирует
	t.Run("closes with two channels, second closes", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		close(ch2)
		orDone := Or(ch1, ch2)

		select {
		case <-orDone:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Expected orDone to be closed")
		}
	})

	//Проверка что на постоянных вызовах и гонках он не отпадает
	t.Run("stress test with goroutines", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ch1 := make(chan interface{})
			ch2 := make(chan interface{})
			orDone := Or(ch1, ch2)

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				close(ch1)
			}()

			go func() {
				defer wg.Done()
				<-orDone
			}()

			wg.Wait()
		}
	})

	//Проверить функция устойчива к гонкам при одновременном закрытии нескольких каналов
	t.Run("race condition test", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			ch1 := make(chan interface{})
			ch2 := make(chan interface{})
			ch3 := make(chan interface{})

			orDone := Or(ch1, ch2, ch3)

			go func() { close(ch1) }()
			go func() { close(ch2) }()
			go func() { close(ch3) }()

			<-orDone
		}
	})

	//Проверяет что не падает при передаче уже закрытых каналов
	t.Run("nil channels are not allowed as input", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		close(ch1)
		orDone := Or(ch1, ch2)

		select {
		case <-orDone:
		case <-time.After(10 * time.Millisecond):
			t.Fatal("Expected orDone to be closed")
		}
	})
}
