package kraken

import (
	"github.com/jfixby/kraken/input"
	"github.com/jfixby/kraken/orderbook"
	testoutput "github.com/jfixby/kraken/output"
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/fileops"
	"path/filepath"
	"testing"
	"time"
)

var setup *testing.T

func TestOrderbook(t *testing.T) {
	setup = t

	home := fileops.Abs("")
	testData := filepath.Join(home, "data", "test1")
	testOutput := filepath.Join(testData, "out", "output_file.csv")
	testInput := filepath.Join(testData, "in", "input_file.csv")

	test := &testoutput.TestOutput{File: testOutput}
	test.LoadAll()

	reader := input.NewFileReader(testInput)
	testListener := &TestListener{
		testData: test}
	reader.Subscribe(testListener)
	reader.Run()

	var bookEventListener orderbook.BookListener = testListener
	book := orderbook.NewBook(bookEventListener)
	testListener.book = book

	for reader.IsRunnung() {
		time.Sleep(2 * time.Second)
	}

	pin.D("EXIT")
}

type TestListener struct {
	testData *testoutput.TestOutput
	scenario string
	book     *orderbook.Book
	counter  int
}

func (t *TestListener) DoProcess(ev *orderbook.Event) {
	t.book.DoUpdate(ev)
	pin.D("Event received", ev)
}

func (t *TestListener) OnBookEvent(e *orderbook.BookEvent) {
	expectedEvent := t.testData.GetEvent(t.scenario, t.counter)
	check(setup, e, expectedEvent, t.scenario, t.counter)
	t.counter++
}

func check(
	setup *testing.T,
	actual *orderbook.BookEvent,
	expected *orderbook.BookEvent,
	scenario string,
	counter int) {

	if !expected.Equal(actual) {

		pin.D(" counter", counter)
		pin.D("expected", expected)
		pin.D("  actual", actual)
		//setup.FailNow()
		panic("")
	}
}

func (t *TestListener) Reset(scenario string) {
	pin.D("Next scenario", scenario)
	t.scenario = scenario
	t.counter = 0
}