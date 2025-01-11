package main

import "strings"

// how do I get the uses of my API to actually use my builders
// as opposed to stop messing with the objects directly?
// And one approach to this is you simply hide the objects that you want your users not to touch.

type email struct {
	from, to, subject, body string
}

type EmailBuilder struct {
	email email
}

func (b *EmailBuilder) From(from string) *EmailBuilder {
	if !strings.Contains(from, "@") {
		panic("email should contain @")
	}

	b.email.from = from
	return b
}

func (b *EmailBuilder) To(to string) *EmailBuilder {
	b.email.to = to
	return b
}

func (b *EmailBuilder) Subject(subject string) *EmailBuilder {
	b.email.subject = subject
	return b
}

func (b *EmailBuilder) Body(body string) *EmailBuilder {
	b.email.body = body
	return b
}

func sendMailImpl(email *email) {

}

// you don't want your clients to actually work with the email object.
// You only want to work with a builder. So how do you do this?
// Well, you can do this by using a builder parameter, and that's basically going to be a
// function which sort of applies to the builder. So you have to provide a function which takes
// a builder and then does something with it, typically sort of calls, something on the builder....

type build func(*EmailBuilder)

func SendEmail(action build) {
	builder := EmailBuilder{}

	action(&builder)
	//  +---------------------------------------------------------------------------------------------------------+
	//  | so here after the action function all the necessary stuff will be initialized                           |
	//  | by this passed function, without leaking important stuff that may be missed with by clients             |
	//  | and the hidden sendMailImpl function will have all the required parameters to process sending the email |
	//  +---------------------------------------------------------------------------------------------------------+
	sendMailImpl(&builder.email)
}

func main() {
	SendEmail(func(b *EmailBuilder) {
		//  +-------------------------------------------------------------------------------+
		//  | 		                                                                           |
		//  |     And this is precisely the location where you don't have access to the     |
		//  | 		email object itself. You only have access to the builder and you can use   |
		//  | 		that builder to build up information about the email.                      |
		//  | 		                                                                           |
		//  +-------------------------------------------------------------------------------+

		b.From("foo@bar.com").
			To("bar@baz.com").
			Subject("Meeting").
			Body("Hello, do you want to meet?")
	})

	//	So this is yet another approach to how you can use the builder, how in fact you can force clients
	//
	// to use the builder as opposed to providing some sort of incomplete object for initialization,
}
