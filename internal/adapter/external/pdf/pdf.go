package pdf

import (
	"bytes"
	"log"
	"github.com/ledongthuc/pdf"
)

type PDFService interface {
	ReadPDF(fileName string) (string, error)
}

type PdfServiceImpl struct{}

func NewPDFService() *PdfServiceImpl {
	return &PdfServiceImpl{}
}

func (p *PdfServiceImpl) ReadPDF(fileName string) (string, error) {
	f, r, err := pdf.Open(fileName)

	if err != nil {
		log.Println("Error in opening pdf", err)
		return "", err
	}
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		log.Println("Error in reading pdf", err)
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
