package games

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/internal"
)

func (svc *service) CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	if pdf.Name == "" {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeInvalidPDFName,
			Detail: "The given name cannot be empty.",
		})
	}

	if len(pdf.Name) > 30 {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeInvalidPDFName,
			Detail: "The given name cannot go over 30 characters.",
		})
	}

	pdf.OwnerID = session.UserID
	pdf.Pages = make([]map[string]string, len(internal.PDFSchemaPageNames[pdf.Schema]))

	err := svc.schema.PDFInsert(ctx, pdf)
	if err != nil {
		return errors.Wrap(err, "cannot insert PDF")
	}

	if pdf.GameID != nil {
		subject := fmt.Sprintf("games.%s", pdf.GameID)
		err = svc.PubSub.Publish(ctx, subject, &internal.EventWrapper[any]{
			Event:   internal.EventPDF,
			Payload: pdf,
		})
	}

	return errors.Wrap(err, "cannot publish PDF event")
}

func (svc *service) GetPDF(ctx context.Context, pdfID uuid.UUID, view internal.PDFView) (*internal.PDF, error) {
	pdf, err := svc.schema.PDFGet(ctx, pdfID, view)
	return pdf, errors.Wrap(err, "cannot get PDF")
}

func (svc *service) GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error) {
	if pageNum < 1 {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeInvalidPDFPageNumber,
		})
	}

	pageFields, err := svc.schema.PDFGetPage(ctx, pdfID, pageNum)
	return pageFields, errors.Wrap(err, "cannot get PDF page")
}

func (svc *service) UpdatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	err := svc.schema.PDFUpdate(ctx, session, pdf)
	return errors.Wrap(err, "cannot update PDF")
}

func (svc *service) UpdatePDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string, broadcast bool) error {
	if pageNum < 1 {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeInvalidPDFPageNumber,
		})
	}

	if fieldName == "" {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeInvalidPDFFieldName,
			Detail: "The given field name cannot be empty.",
		})
	}

	err := svc.schema.PDFUpdatePage(ctx, pdfID, pageNum, fieldName, fieldValue)
	if err != nil {
		return errors.Wrap(err, "cannot update PDF page")
	}

	if broadcast {
		subject := fmt.Sprintf("pdfs.%s.pages.%v", pdfID, pageNum)
		err = svc.PubSub.Publish(ctx, subject, &internal.EventWrapper[any]{
			Event: internal.EventPDFPage,
			Payload: &internal.PDFPage{
				PDFID:   pdfID,
				PageNum: pageNum,
				Fields: map[string]string{
					fieldName: fieldValue,
				},
			},
		})
	}

	return errors.Wrap(err, "cannot publish PDF_PAGE event")
}

func (svc *service) DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID, gameID *uuid.UUID) error {
	err := svc.schema.PDFDelete(ctx, pdfID, session.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot delete PDF")
	}

	if gameID != nil {
		subject := fmt.Sprintf("games.%s", gameID)
		err = svc.PubSub.Publish(ctx, subject, &internal.EventWrapper[any]{
			Event: internal.EventDeletedPDF,
			Payload: &internal.PDF{
				ID: pdfID,
			},
		})
	}

	return errors.Wrap(err, "cannot publish DELETED_PDF event")
}
