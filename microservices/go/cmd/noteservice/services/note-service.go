package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/mappers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	cDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	kDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	nDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/notedtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
)

type NoteService interface {
	AddNoteTxn(
		ctx context.Context,
		userBo userbos.UserBo,
		dto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
	) (cDTOs.SuccessDto, error)
	UpdateNoteTxn(
		ctx context.Context,
		userBo userbos.UserBo,
		dto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
	) (cDTOs.SuccessDto, error)
	DeleteNoteTxn(
		ctx context.Context,
		userBo userbos.UserBo,
		noteIdDto nDTOs.NoteIdDto,
	) (cDTOs.SuccessDto, error)
	GetNoteById(
		ctx context.Context,
		userBo userbos.UserBo,
		sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto],
	) (nDTOs.NoteReadDto, error)
	GetNotesPage(
		ctx context.Context,
		userBo userbos.UserBo,
		sessReqDto cDTOs.UKeySessionReqDto[pagination.PageRequest],
	) (pagination.Page[nDTOs.NotePreviewDto], error)
	DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error)
}

type NoteServiceImpl struct {
	noteRepository repositories.NoteRepository
	userKeyService externalservices.ExtUserKeyService
	crudDSHandler  dshandlers.CrudDSHandler
	errorService   sharedservices.ErrorService
	noteBr         businessrules.NoteBr
}

func (n NoteServiceImpl) AddNoteTxn(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
) (cDTOs.SuccessDto, error) {
	return dshandlers.Txn(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) (cDTOs.SuccessDto, error) {
			return n.AddNoteTxn(ctx, userBo, sessReqDto)
		})
}

func (n NoteServiceImpl) addNote(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
) (cDTOs.SuccessDto, error) {
	sessDto, noteCreateDto := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
	keyDto, err := n.userKeyService.GetKeyFromSession(ctx, sessDto)
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	key, err := kDTOs.UserKeyDto.GetKey(keyDto)
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	titleCipher, err := cipherutils.EncryptAES(key, []byte(noteCreateDto.Title))
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	textCipher, err := cipherutils.EncryptAES(key, []byte(noteCreateDto.Text))
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	note := models.Note{
		UserId:      userBo.Id,
		TextCipher:  textCipher,
		TitleCipher: titleCipher,
		KeyVersion:  keyDto.KeyVersion,
	}
	if _, err := n.noteRepository.Create(ctx, note); err != nil {
		return cDTOs.SuccessDto{}, err
	}
	return cDTOs.NewSuccessTrue(), nil
}

func (n NoteServiceImpl) UpdateNoteTxn(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
) (cDTOs.SuccessDto, error) {
	return dshandlers.Txn(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) (cDTOs.SuccessDto, error) {
			return n.updateNote(ctx, userBo, sessReqDto)
		})
}

func (n NoteServiceImpl) updateNote(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
) (cDTOs.SuccessDto, error) {
	sessDto, noteUpdateDto := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
	existingNote, err := n.getExistingNote(ctx, noteUpdateDto.Id)
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	keyDto, err := n.userKeyService.GetKeyFromSession(ctx, sessDto)
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	if err := n.noteBr.ValidateNoteUpdate(userBo, keyDto, existingNote); err != nil {
		return cDTOs.SuccessDto{}, err
	}
	key, err := keyDto.GetKey()
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	titleCipher, err := cipherutils.EncryptAES(key, []byte(noteUpdateDto.Title))
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	textCipher, err := cipherutils.EncryptAES(key, []byte(noteUpdateDto.Text))
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	existingNote.TitleCipher = titleCipher
	existingNote.TextCipher = textCipher
	existingNote.KeyVersion = keyDto.KeyVersion
	if _, err := n.noteRepository.Update(ctx, existingNote); err != nil {
		return cDTOs.SuccessDto{}, err
	}
	return cDTOs.NewSuccessTrue(), nil
}

func (n NoteServiceImpl) DeleteNoteTxn(
	ctx context.Context,
	userBo userbos.UserBo,
	noteIdDto nDTOs.NoteIdDto,
) (cDTOs.SuccessDto, error) {
	return dshandlers.Txn(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) (cDTOs.SuccessDto, error) {
			return n.deleteNote(ctx, userBo, noteIdDto)
		})
}

func (n NoteServiceImpl) deleteNote(
	ctx context.Context,
	userBo userbos.UserBo,
	noteIdDto nDTOs.NoteIdDto,
) (cDTOs.SuccessDto, error) {
	existingNote, err := n.getExistingNote(ctx, noteIdDto.Id)
	if err != nil {
		return cDTOs.SuccessDto{}, err
	}
	if err := n.noteBr.ValidateNoteDelete(userBo, existingNote); err != nil {
		return cDTOs.SuccessDto{}, err
	}
	if _, err := n.noteRepository.Delete(ctx, existingNote); err != nil {
		return cDTOs.SuccessDto{}, err
	}
	return cDTOs.NewSuccessTrue(), nil
}

func (n NoteServiceImpl) GetNoteById(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto],
) (nDTOs.NoteReadDto, error) {
	sessDto, noteIdDto := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
	existingNote, err := n.getExistingNote(ctx, noteIdDto.Id)
	if err != nil {
		return nDTOs.NoteReadDto{}, err
	}
	keyDto, err := n.userKeyService.GetKeyFromSession(ctx, sessDto)
	if err != nil {
		return nDTOs.NoteReadDto{}, err
	}
	if err := n.noteBr.ValidateNoteRead(userBo, keyDto, existingNote); err != nil {
		return nDTOs.NoteReadDto{}, err
	}
	key, err := kDTOs.UserKeyDto.GetKey(keyDto)
	if err != nil {
		return nDTOs.NoteReadDto{}, err
	}
	textBytes, err := cipherutils.DecryptAES(key, existingNote.TextCipher)
	if err != nil {
		return nDTOs.NoteReadDto{}, err
	}
	text := string(textBytes)
	titleBytes, err := cipherutils.DecryptAES(key, existingNote.TitleCipher)
	if err != nil {
		return nDTOs.NoteReadDto{}, err
	}
	title := string(titleBytes)

	coreNoteDetails := nDTOs.NewCoreNoteDetailsDto(title, text)
	noteDetailsDto := nDTOs.NoteReadDto{}
	mappers.MapCoreNoteDetailsAndNoteToNoteReadDto(&coreNoteDetails, &existingNote, &noteDetailsDto)
	return noteDetailsDto, nil
}

func (n NoteServiceImpl) GetNotesPage(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[pagination.PageRequest],
) (pagination.Page[nDTOs.NotePreviewDto], error) {
	sessionDto, pageRequest := sessReqDto.SetUserIdAndUnwrap(userBo.Id)

	if err := n.noteBr.ValidateGetNotes(pageRequest); err != nil {
		return pagination.Page[nDTOs.NotePreviewDto]{}, err
	}

	keyDto, err := n.userKeyService.GetKeyFromSession(ctx, sessionDto)
	if err != nil {
		return pagination.Page[nDTOs.NotePreviewDto]{}, err
	}
	key, err := kDTOs.UserKeyDto.GetKey(keyDto)
	if err != nil {
		return pagination.Page[nDTOs.NotePreviewDto]{}, err
	}

	count, err := n.noteRepository.CountByUserId(ctx, userBo.Id)
	if err != nil {
		return pagination.Page[nDTOs.NotePreviewDto]{}, err
	}

	notes, err := n.noteRepository.GetPaginatedByUserId(ctx, userBo.Id, pageRequest)
	if err != nil {
		return pagination.Page[nDTOs.NotePreviewDto]{}, err
	}

	noteDTOs := make([]nDTOs.NotePreviewDto, 0, len(notes))
	for _, note := range notes {
		if note.KeyVersion != keyDto.KeyVersion {
			ruleErr := n.errorService.RuleErrorFromCode(apperrors.ErrCodeDataRace)
			return pagination.Page[nDTOs.NotePreviewDto]{}, apperrors.NewBadReqErrorFromRuleError(ruleErr)
		}
		txtBytes, err := cipherutils.DecryptAES(key, note.TextCipher)
		if err != nil {
			return pagination.Page[nDTOs.NotePreviewDto]{}, err
		}
		titleBytes, err := cipherutils.DecryptAES(key, note.TitleCipher)
		if err != nil {
			return pagination.Page[nDTOs.NotePreviewDto]{}, err
		}
		txt, title := string(txtBytes), string(titleBytes)
		textPreview := utils.StringFirstNChars(txt, 60)
		coreNoteDto := nDTOs.NewCoreNoteDto(title)
		noteReadDto := nDTOs.NotePreviewDto{}
		mappers.MapTextPreviewAndCoreNoteAndNoteToNotePreviewDto(textPreview, &coreNoteDto, &note, &noteReadDto)
		noteDTOs = append(noteDTOs, noteReadDto)
	}
	return pagination.NewPage(noteDTOs, count), nil
}

func (u NoteServiceImpl) DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error) {
	return u.noteRepository.DeleteByUserIdAndGetCount(ctx, userId)
}

func (n NoteServiceImpl) getExistingNote(ctx context.Context, id string) (models.Note, error) {
	noteSearch, err := n.noteRepository.FindById(ctx, id)
	if err != nil {
		return models.Note{}, err
	}
	if note, ok := noteSearch.Get(); ok {
		return note, nil
	} else {
		ruleErr := n.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
		return models.Note{}, apperrors.NewBadReqErrorFromRuleError(ruleErr)
	}
}

func NewNoteServiceImpl(
	noteRepository repositories.NoteRepository,
	userKeyService externalservices.ExtUserKeyService,
	crudDSHandler dshandlers.CrudDSHandler,
	errorService sharedservices.ErrorService,
	noteBr businessrules.NoteBr,
) *NoteServiceImpl {
	return &NoteServiceImpl{
		noteRepository: noteRepository,
		userKeyService: userKeyService,
		crudDSHandler:  crudDSHandler,
		errorService:   errorService,
		noteBr:         noteBr,
	}
}
