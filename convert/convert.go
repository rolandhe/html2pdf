package convert

import (
	"encoding/base64"
	"log"
)

func ConvertPdf(req *PdfRequest,pid int) *ResultML[string] {
	log.Printf("pid=%d,body: %s\n", pid,req.Body)
	task := SendPdfTask(req.Body)
	task.Wait()

	log.Println("build pdf finish")
	err := task.Err
	if err != nil {
		log.Println(err)
		return FailByErrorML[string](err)
	}
	result := base64.StdEncoding.EncodeToString(task.Out)
	//log.Println(result)
	return SuccessML(result)
}
