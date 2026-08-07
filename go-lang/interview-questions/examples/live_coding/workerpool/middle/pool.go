// Middle-версия: те же задачи, но пул пригоден к эксплуатации.
// Разбор отличий от junior — в ../../../../answers.md, раздел «Сделать воркер-пул».
//
// Три идеи, на которых все держится:
//  1. канал tasks НЕ закрывается никогда — значит send on closed невозможен
//     в принципе; воркеры выходят по сигнальным каналам stop/done;
//  2. две фазы остановки: stop = «новые задачи не принимаем, доедаем очередь»,
//     done = «бросаем все, дедлайн вышел»;
//  3. задача — func(ctx) error: ее можно отменить и у нее есть результат.
//
// Что осталось за бортом (появится в senior): выбор стратегии при полной
// очереди и изменение числа воркеров на ходу.
package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// Pool - worker pool с graceful shutdown, panic recovery, per-task context
// и неблокирующим возвратом ошибок.
type Pool struct {
	// Ни одного мьютекса вокруг канала: состояние пула выражено самими каналами.
	// sync.Once нужен только чтобы close(stop) случился ровно один раз.
	tasks chan task
	errs  chan error
	stop  chan struct{} // закрывается в начале Shutdown
	done  chan struct{} // закрывается на force-stop

	wg            sync.WaitGroup
	once          sync.Once
	droppedErrors atomic.Uint64
	// НЮАНС: ошибки, не влезшие в errs, не теряются молча — их считают.
	// Без мониторинга этого счетчика ErrBuf легко оказывается мал.
}

type task struct {
	ctx context.Context //nolint:containedctx // ctx живёт ровно одну задачу, как у http.Request.ctx
	fn  func(ctx context.Context) error
}

// New валидирует cfg и стартует cfg.Workers горутин.
// Паника на невалидном cfg - программерская ошибка.
func New(cfg Config) *Pool {
	cfg.validate()

	p := &Pool{
		tasks: make(chan task, cfg.QueueSize),
		errs:  make(chan error, cfg.ErrBuf),
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}

	p.wg.Add(cfg.Workers)
	// Add на все воркеры сразу и ДО первого go — Wait не может проскочить.
	for range cfg.Workers {
		go p.worker()
	}
	return p
}

// Submit отправляет задачу в очередь.
//
// Возвращает:
//   - nil           - задача принята;
//   - ctx.Err()     - ctx отменён до того, как задача попала в очередь;
//   - ErrPoolClosed - Shutdown уже вызван;
//   - ErrNilTask    - fn == nil.
func (p *Pool) Submit(ctx context.Context, fn func(ctx context.Context) error) error {
	if fn == nil {
		return ErrNilTask
	}

	// Fast-path: неблокирующе проверяем «пул закрывается» и «ctx отменен».
	// Без него select ниже при нескольких готовых ветвях выбрал бы случайную
	// и мог протолкнуть задачу в уже закрывающийся пул.
	select {
	case <-p.stop:
		return ErrPoolClosed
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	t := task{ctx: ctx, fn: fn}
	// ИЗВЕСТНОЕ ОГРАНИЧЕНИЕ middle: если close(stop) случится, когда мы уже
	// внутри этого select, готовы обе ветви и выбор случайный — задача может
	// лечь в буфер после начала Shutdown. Обычно ее подхватит drainAndExit,
	// но при force-stop по дедлайну она потеряется молча. Закрывается только
	// линеаризацией Submit против Shutdown (refcount in-flight отправителей),
	// то есть синхронизацией на горячем пути — здесь сознательно не сделано.
	select {
	case p.tasks <- t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.stop:
		return ErrPoolClosed
	}
}

// Shutdown - graceful drain с force-stop fallback'ом.
//
// Возвращает:
//   - nil                 - drain успел до дедлайна;
//   - обёрнутый ctx.Err() - ctx истёк (errors.Is с DeadlineExceeded);
//   - ErrAlreadyShutdown  - повторный вызов.
func (p *Pool) Shutdown(ctx context.Context) error {
	called := false
	p.once.Do(func() {
		called = true
		close(p.stop)
	})
	if !called {
		return ErrAlreadyShutdown
	}

	drained := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(drained)
	}()

	select {
	case <-drained:
		return nil
	case <-ctx.Done():
		// Дедлайн вышел → force-stop. Хвост очереди дропается молча:
		// в проде это место логируют и снимают метрику len(p.tasks).
		close(p.done)
		<-drained
		return fmt.Errorf("workerpool: дедлайн на drain истёк: %w", ctx.Err())
	}
}

// Errors - канал ошибок и паник от воркеров. НЕ закрывается на Shutdown:
// in-flight задачи могут писать после возврата Shutdown. Паники приходят
// как *PanicError.
func (p *Pool) Errors() <-chan error {
	// НЮАНС: канал ошибок не закрывается принципиально — задача, уже начавшая
	// выполняться, может записать в него после возврата Shutdown, а запись
	// в закрытый канал паникует. Потребитель выходит по своему сигналу.
	return p.errs
}

// DroppedErrors - счётчик ошибок, не вместившихся в Errors().
func (p *Pool) DroppedErrors() uint64 {
	return p.droppedErrors.Load()
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.done:
			return
		case t := <-p.tasks:
			p.runTask(t)
		case <-p.stop:
			// Не выходим сразу: сначала дочерпываем очередь, иначе Shutdown
			// вернул бы nil, оставив невыполненные задачи в буфере.
			p.drainAndExit()
			return
		}
	}
}

func (p *Pool) drainAndExit() {
	for {
		select {
		case <-p.done:
			return
		case t := <-p.tasks:
			p.runTask(t)
		default:
			return
		}
	}
}

// runTask - отдельная функция, чтобы defer recover
// срабатывал на каждой задаче, а не один раз на всю жизнь воркера.
func (p *Pool) runTask(t task) {
	// recover ставится на КАЖДУЮ задачу (отдельная функция), а не один раз
	// на жизнь воркера: упавшая задача не должна уносить с собой воркер.
	defer func() {
		if r := recover(); r != nil {
			p.sendErr(&PanicError{Recovered: r, Stack: debug.Stack()})
		}
	}()

	if err := t.ctx.Err(); err != nil {
		// Задача могла пролежать в очереди дольше своего дедлайна — проверяем
		// перед запуском, чтобы не делать заведомо ненужную работу.
		p.sendErr(fmt.Errorf("workerpool: ctx задачи отменён до запуска: %w", err))
		return
	}

	if err := t.fn(t.ctx); err != nil {
		p.sendErr(err)
	}
}

func (p *Pool) sendErr(err error) {
	// Отправка ошибки неблокирующая: медленный потребитель не должен
	// останавливать воркеров. Цена — потерянные ошибки, они в счетчике.
	select {
	case p.errs <- err:
	default:
		p.droppedErrors.Add(1)
	}
}
