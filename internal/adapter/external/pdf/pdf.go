package pdf

import (
	"log"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func (p *PDF) ConvertToPDF(fileName string, data [][]string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)

	pdf.AddPage()

	lineHeight := 10.0
	marginLeft := 10.0
	pageWidth := 190.0 // A4 width minus margins
	columnSpacing := 5.0
	columnWidths := calculateColumnWidths(data, pageWidth)

	for _, row := range data {
		pdf.SetX(marginLeft)

		maxRowHeight := lineHeight // Track the maximum height of the current row

		for i, col := range row {
			// Sanitize the string
			col = strings.TrimSpace(col)

			// Define the width for each column
			colWidth := columnWidths[i]

			// Store the current X position
			startX := pdf.GetX()

			// Add the text with wrapping for each column
			pdf.MultiCell(colWidth, lineHeight, col, "", "L", false)

			// Update max row height based on the current MultiCell height
			rowHeight := pdf.GetY() - pdf.GetY()
			if rowHeight > maxRowHeight {
				maxRowHeight = rowHeight
			}

			// Reset X and move to the next column
			pdf.SetXY(startX+colWidth+columnSpacing, pdf.GetY())
		}

		// Move to the next row
		pdf.Ln(maxRowHeight)
	}

	// Save the PDF
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting the current working directory:", err)
		return err
	}

	uploadDir := dir + "/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			log.Println("Error creating uploads directory:", err)
			return err
		}
	}

	filePath := uploadDir + "/" + fileName + ".pdf"
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Printf("Error saving PDF: %v", err)
		return err
	}

	log.Println("PDF generated successfully at", filePath)
	return nil
}

// calculateColumnWidths determines the width of each column based on the content
func calculateColumnWidths(data [][]string, totalWidth float64) []float64 {
	if len(data) == 0 {
		return []float64{}
	}

	numCols := len(data[0])
	defaultWidth := totalWidth / float64(numCols)
	widths := make([]float64, numCols)

	// Calculate column widths based on the maximum length of content in each column
	for col := 0; col < numCols; col++ {
		maxLength := 0
		for _, row := range data {
			if col < len(row) && len(row[col]) > maxLength {
				maxLength = len(row[col])
			}
		}

		// Adjust width proportionally (tweak factor as needed)
		widths[col] = defaultWidth + float64(maxLength)*0.5
	}

	// Ensure total width does not exceed page width
	totalCalculatedWidth := 0.0
	for _, width := range widths {
		totalCalculatedWidth += width
	}
	if totalCalculatedWidth > totalWidth {
		scaleFactor := totalWidth / totalCalculatedWidth
		for i := range widths {
			widths[i] *= scaleFactor
		}
	}

	return widths
}

