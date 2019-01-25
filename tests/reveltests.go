package tests

import (
	"github.com/revel/revel/testing"
)

type RevelTest struct {
	testing.TestSuite
}

func (t *RevelTest) Before() {
	println("Set up")
}

func (t *RevelTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *RevelTest) After() {
	println("Tear down")
}
