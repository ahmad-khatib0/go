package main

import "fmt"

// a single builder is sufficient for building up a particular object. But there are
// situations where you need more than one builder, where you need to somehow separate the
// process of building up the different aspects of a particular type.

type Person struct {
	StreetAddress, Postcode, City string
	CompanyName, Position         string
	AnnualIncome                  int
}

type PersonBuilder struct {
	person *Person
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{&Person{}}
}

// &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
//  ┌───────────────────────────────────────────────────────────────────────────────────────────────────┐
//    So now we have ways of transitioning from a  PersonBuilder to either a person
//    address builder and a person job builder. But you have to realize that effectively:
//    person, job builder and person address builder are both person builders. And as a result,
//    when you have a person address builder, you can quickly use the works method to jump to a person,
//    job builder and vice versa. You can jump back to the person address builder using the lives method.
//    So that's very convenient.
//  └───────────────────────────────────────────────────────────────────────────────────────────────────┘
// &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&

type PersonAddressBuilder struct {
	PersonBuilder // aggregates PersonBuilder
}

func (b *PersonBuilder) Lives() *PersonAddressBuilder {
	return &PersonAddressBuilder{*b}
}

func (it *PersonAddressBuilder) At(streetAddress string) *PersonAddressBuilder {
	it.person.StreetAddress = streetAddress
	return it
}

func (it *PersonAddressBuilder) In(city string) *PersonAddressBuilder {
	it.person.City = city
	return it
}

func (it *PersonAddressBuilder) WithPostcode(postcode string) *PersonAddressBuilder {
	it.person.Postcode = postcode
	return it
}

type PersonJobBuilder struct {
	PersonBuilder // aggregates PersonBuilder
}

func (b *PersonBuilder) Works() *PersonJobBuilder {
	return &PersonJobBuilder{*b}
}

func (pjb *PersonJobBuilder) At(companyName string) *PersonJobBuilder {
	pjb.person.CompanyName = companyName
	return pjb
}

func (pjb *PersonJobBuilder) AsA(position string) *PersonJobBuilder {
	pjb.person.Position = position
	return pjb
}

func (pjb *PersonJobBuilder) Earning(annualIncome int) *PersonJobBuilder {
	pjb.person.AnnualIncome = annualIncome
	return pjb
}

func (b *PersonBuilder) Build() *Person {
	return b.person
}

func main() {

	pb := NewPersonBuilder()
	pb.
		Lives().
		At("123 London Road").
		In("London").
		WithPostcode("SW12BC").
		// switch from one builder to a completely different builder.
		Works().
		At("Fabrikam").
		AsA("Programmer").
		Earning(123000)

	person := pb.Build()

	fmt.Println(*person)
	// {123 London Road SW12BC London Fabrikam Programmer 123000}
}
