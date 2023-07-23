package main

// The interface segregation principle is a really simple principle, it's probably the simplest principle
// out of the solid design principles. Basically, what it states is that you shouldn't put too much into
// an interface. You shouldn't try to throw everything and the kitchen sink into just one single interface.
// And sometimes it makes sense to break up the interface into several smaller interfaces.

type Document struct {
}

type Machine interface {
	Print(doc Document)
	Fax(doc Document)
	Scan(doc Document)
}

// ok if you need a multifunction device
type MultiFunctionPrinter struct {
	// ...
}

func (m MultiFunctionPrinter) Print(d Document) {}
func (m MultiFunctionPrinter) Fax(d Document)   {}
func (m MultiFunctionPrinter) Scan(d Document)  {}

// An old fashioned printer doesn't really have any scanning or faxing capabilities.
// But because you want to implement this interface, because maybe some other APIs rely on the machine
// interface, you have to implement this anyway. You are being forced into implementing it.
// So you go ahead through all the similar motions to implement the machine interface and you end up with
// the same stuff as you would for a multifunction device, except there is a bit of a problem.
// the vast majority of them is not going to work with this OldFashionedPrinter
type OldFashionedPrinter struct{}

func (o OldFashionedPrinter) Print(d Document) {
	// ok
}

// not a good solution
func (o OldFashionedPrinter) Fax(d Document) {
	panic("operation not supported")
}

// Deprecated: ...
// not a good solution, because maybe some IDEs doesn't  show a deprecation message
func (o OldFashionedPrinter) Scan(d Document) {
	panic("operation not supported")
}

// #################################
// ISP
// So the interface aggregation principle basically states that try to break up an
// interface into separate parts that people will definitely need. So there is no
// guarantee that if somebody needs printing, they also need faxing. So in this particular
// example, it might make more sense to split up the printing and the scanning into separate interfaces.

// better approach: split into several interfaces
type Printer interface {
	Print(d Document)
}

type Scanner interface {
	Scan(d Document)
}

// so if you want to print Documents only, you then implement the Printer interface

// printer only
type MyPrinter struct{}

func (m MyPrinter) Print(d Document) {}

// combine interfaces
type Photocopier struct{}

func (p Photocopier) Scan(d Document)  {}
func (p Photocopier) Print(d Document) {}

type MultiFunctionDevice interface {
	Printer
	Scanner
}

// interface combination + decorator
type MultiFunctionMachine struct {
	printer Printer
	scanner Scanner
}

func (m MultiFunctionMachine) Print(d Document) {
	m.printer.Print(d)
}

func (m MultiFunctionMachine) Scan(d Document) {
	m.scanner.Scan(d)
}

//  +-------------------------------------------------------------------------------------------------+
//  | So you can see that with the interface aggregation approach, what you can do is, first of all,  |
//  | you have very granular kind of definitions. So you just grab the interfaces that you need and   |
//  | you don't have any extra members in those interfaces. So if you're just building an ordinary    |
//  | printer, you just get the print method, then that's pretty much it.                             |
//  +-------------------------------------------------------------------------------------------------+

func main() {

}
