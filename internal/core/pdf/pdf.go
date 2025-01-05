package pdf

type PDFService interface {
	ReadPDF(fileName string) error
}
