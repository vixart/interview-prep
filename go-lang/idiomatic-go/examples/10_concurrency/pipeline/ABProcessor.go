// Этап AB конвейера: две независимые задачи в отдельных горутинах,
// результаты и ошибки — в буферизованные каналы, wait собирает оба результата
// через select и умеет прерваться по ctx.Done().
package main

import "context"

type abProcessor struct {
	outA chan aOut
	outB chan bOut
	errs chan error
}

func newABProcessor() *abProcessor {
	return &abProcessor{
		outA: make(chan aOut, 1),
		// буфер 1: горутина отдаст результат и завершится, даже если ее уже никто не ждет
		outB: make(chan bOut, 1),
		errs: make(chan error, 2),
		// место под ошибку от каждой из двух горутин
	}
}

func (p *abProcessor) start(ctx context.Context, data Input) {
	go func() {
		aOut, err := getResultA(ctx, data.A)
		if err != nil {
			p.errs <- err
			return
		}
		p.outA <- aOut
	}()
	go func() {
		bOut, err := getResultB(ctx, data.B)
		if err != nil {
			p.errs <- err
			return
		}
		p.outB <- bOut
	}()
}

func (p *abProcessor) wait(ctx context.Context) (cIn, error) {
	var cData cIn
	for count := 0; count < 2; count++ {
		// ждем ровно два результата, порядок неважен
		select {
		case a := <-p.outA:
			cData.a = a
		case b := <-p.outB:
			cData.b = b
		case err := <-p.errs:
			return cIn{}, err
		case <-ctx.Done():
			return cIn{}, ctx.Err()
		}
	}
	return cData, nil
}

type aOut struct {
}

type bOut struct {
}

type cIn struct {
	a aOut
	b bOut
}

func getResultA(ctx context.Context, in string) (aOut, error) {
	return aOut{}, nil
}

func getResultB(ctx context.Context, in string) (bOut, error) {
	return bOut{}, nil
}
