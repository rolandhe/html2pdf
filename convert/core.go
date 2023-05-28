package convert

import (
	"bytes"
	"errors"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"log"
	"runtime"
)

var (
	taskMQ     = make(chan *PdfTask, 5)
	TimeoutErr = errors.New("get pdfconv engine timeout")
)

type PdfTask struct {
	body       string
	Out        []byte
	Err        error
	notifyChan chan struct{}
}

func SendPdfTask(body string) *PdfTask {
	task := &PdfTask{
		body:       body,
		notifyChan: make(chan struct{}),
	}
	taskMQ <- task
	return task
}

func (task *PdfTask) notify() {
	close(task.notifyChan)
}

func (task *PdfTask) Wait() {
	<-task.notifyChan
}

func InitPdf() error {

	go consumer()
	return nil
}

func consumer() {
	log.Printf("start consumer")
	runtime.LockOSThread()
	var err error
	err = pdf.Init()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		task := <-taskMQ
		log.Println("got task")
		task.Out, task.Err = convertHtml2Pdf(task.body)
		log.Println("finish task")
		task.notify()
	}
}

func convertHtml2Pdf(body string) ([]byte, error) {
	var buff bytes.Buffer
	_, err := buff.WriteString(body)
	if err != nil {
		return nil, err
	}
	return convertPdfCore(&buff)
}

func convertPdfCore(buff *bytes.Buffer) ([]byte, error) {
	var err error

	obj, err := pdf.NewObjectFromReader(buff)
	if err != nil {
		return nil, err
	}
	converter, err := pdf.NewConverter()
	if err != nil {
		return nil, err
	}
	defer converter.Destroy()
	converter.Add(obj)

	// Set converter options.
	converter.Title = "Invoice"
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Portrait
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "10mm"
	converter.MarginRight = "10mm"
	var out bytes.Buffer
	if err = converter.Run(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
