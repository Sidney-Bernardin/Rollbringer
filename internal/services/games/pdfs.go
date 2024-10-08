package games

import (
	"context"
	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *service) CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	pdf.OwnerID = session.UserID
	pdf.Pages = make([]map[string]string, len(internal.PDFSchemaPageNames[pdf.Schema]))

	err := svc.schema.PDFInsert(ctx, pdf)
	return errors.Wrap(err, "cannot insert PDF")
}

func (svc *service) GetPDF(ctx context.Context, pdfID uuid.UUID, view internal.PDFView) (*internal.PDF, error) {
	pdf, err := svc.schema.PDFGet(ctx, pdfID, view)
	return pdf, errors.Wrap(err, "cannot get PDF")
}

func (svc *service) GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error) {
	pageFields, err := svc.schema.PDFGetPage(ctx, pdfID, pageNum)
	return pageFields, errors.Wrap(err, "cannot get PDF page")
}

func (svc *service) UpdatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	err := svc.schema.PDFUpdate(ctx, session, pdf)
	return errors.Wrap(err, "cannot update PDF")
}

func (svc *service) UpdatePDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error {
	err := svc.schema.PDFUpdatePage(ctx, pdfID, pageNum, fieldName, fieldValue)
	return errors.Wrap(err, "cannot update PDF page")
}

func (svc *service) DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID) error {
	err := svc.schema.PDFDelete(ctx, pdfID, session.UserID)
	return errors.Wrap(err, "cannot delete PDF")
}
