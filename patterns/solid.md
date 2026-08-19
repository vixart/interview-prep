# SOLID

Пять принципов проектирования (формулировки — Robert C. Martin).
Разбор с примерами: <https://medium.com/webbdev/solid-4ffc018077da>

## S — Single Responsibility

> A module should be responsible to one, and only one, actor.

У модуля должна быть **одна причина для изменения** — один «актор»
(группа заинтересованных лиц), ради которого его меняют. Не «класс делает
одну вещь», а «за изменения класса отвечает один заказчик»: если отчёт
для бухгалтерии и расчёт зарплаты живут в одном классе, изменения одного
ломают другого.

## O — Open/Closed

> A software artifact should be open for extension but closed for modification.

Поведение расширяется **добавлением нового кода**, а не правкой старого.
Практически: новые варианты поведения — новые реализации интерфейса,
а не новый `case` в разросшемся `switch` по типу.

## L — Liskov Substitution

> If for each object o1 of type S there is an object o2 of type T such that
> for all programs P defined in terms of T, the behavior of P is unchanged
> when o1 is substituted for o2, then S is a subtype of T.

Подтип должен быть **подставим** вместо базового типа без изменения
поведения программы. Реализация интерфейса не имеет права ужесточать
предусловия, ослаблять постусловия или кидать сюрпризы там, где контракт
их не обещал (классика: `Square` наследует `Rectangle` и ломает `SetWidth`).

## I — Interface Segregation

> The lesson here is that depending on something that carries baggage that
> you don't need can cause you troubles that you didn't expect.

Не заставляй клиента зависеть от методов, которые ему не нужны: зависимость
от «жирного» интерфейса тянет за собой чужой багаж — перекомпиляцию,
переусложнение моков, ложную связность. Лучше несколько узких интерфейсов.
В Go это идиома по умолчанию: маленькие интерфейсы (`io.Reader`,
`io.Writer`) объявляются на стороне **потребителя**.

## D — Dependency Inversion

> The most flexible systems are those in which source code dependencies
> refer only to abstractions, not to concretions.

Зависимости в исходном коде должны указывать на **абстракции**, а не на
конкретные реализации. Бизнес-логика объявляет интерфейс «что мне нужно»
(репозиторий, шлюз), инфраструктура его реализует; направление зависимостей —
внутрь, к домену, а не наружу к деталям.

## Для собеседования

- SRP — про причину изменения (актора), а не про «одну функцию».
- OCP достигается через полиморфизм/интерфейсы, а не через флаги и `if`.
- LSP — поведенческий контракт, компилятор его не проверяет.
- ISP в Go — «принимай интерфейсы, возвращай структуры», интерфейс у клиента.
- DIP ≠ dependency injection: DI — техника, DIP — принцип о направлении
  зависимостей (DI — один из способов его соблюсти).
