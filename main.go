package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type InvoiceItem struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type InvoiceRequest struct {
	FreelancerName string        `json:"freelancer_name"`
	ClientName     string        `json:"client_name"`
	InvoiceNumber  string        `json:"invoice_no"`
	Items          []InvoiceItem `json:"items"`
}

func main() {
	r := gin.Default()

	// 1. CORS Middleware - මේක තමයි අර රතු පාට Error එක නිවැරදි කරන්නේ
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/generate-invoice", func(c *gin.Context) {
		var req InvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid Data"})
			return
		}

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()

		// Logo එක (logo.png file එක folder එකේ තිබිය යුතුයි)
		pdf.ImageOptions("logo.png", 10, 10, 30, 0, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")

		// Header Design
		pdf.SetTextColor(0, 102, 204)
		pdf.SetFont("Arial", "B", 20)
		pdf.CellFormat(0, 10, req.FreelancerName, "", 1, "R", false, 0, "")
		pdf.Ln(20)

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(0, 10, fmt.Sprintf("Invoice No: %s", req.InvoiceNumber))
		pdf.Ln(8)
		pdf.Cell(0, 10, fmt.Sprintf("Bill To: %s", req.ClientName))
		pdf.Ln(15)

		// Table Header
		pdf.SetFillColor(0, 102, 204)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(140, 10, "Description", "1", 0, "C", true, 0, "")
		pdf.CellFormat(50, 10, "Price ($)", "1", 1, "C", true, 0, "")

		// Table Body
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 12)
		var total float64
		for _, item := range req.Items {
			pdf.CellFormat(140, 10, item.Description, "1", 0, "L", false, 0, "")
			pdf.CellFormat(50, 10, fmt.Sprintf("%.2f", item.Price), "1", 1, "R", false, 0, "")
			total += item.Price
		}

		// Total
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(140, 10, "Total Amount", "1", 0, "R", false, 0, "")
		pdf.CellFormat(50, 10, fmt.Sprintf("%.2f", total), "1", 1, "R", false, 0, "")

		// 2. Stream PDF - මේකෙන් තමයි browser එකට download එක එවන්නේ
		c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
		c.Header("Content-Type", "application/pdf")
		pdf.Output(c.Writer)
	})

	r.Run(":9090")
}
