package pdf

type PDFService interface {
	ConvertToPDF(fileName string,data [][]string) error
}


type PDF struct{

}

func NewPDFService() *PDF {
	return &PDF{}
}