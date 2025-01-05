package pdf

import (
	"file-modification/internal/adapter/external/pdf"
	"fmt"
	"log"
)

type PdfServiceImpl struct {
	PdfReader pdf.PDFService
}

func NewPdfServiceImpl(pdf pdf.PDFService) *PdfServiceImpl {
	return &PdfServiceImpl{
		PdfReader: pdf,
	}
}

func (p *PdfServiceImpl) ReadPDF(fileName string)( error) {
	fileName="uploads/resume.pdf"

	str,err:=p.PdfReader.ReadPDF(fileName)
	if err!=nil{

		log.Printf("Error in reading the pdf %s.Error is %v\n", fileName, err)
		return  err

	}

	fmt.Println(str)
	return nil
}
