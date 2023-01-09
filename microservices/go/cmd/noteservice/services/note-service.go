package services

import (
	"context"
	"github.com/barweiss/go-tuple"
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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type NoteService interface {
	AddNoteTransaction(
		ctx context.Context,
		userBo userbos.UserBo,
		dto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
	) single.Single[cDTOs.SuccessDto]
	UpdateNoteTransaction(
		ctx context.Context,
		userBo userbos.UserBo,
		dto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
	) single.Single[cDTOs.SuccessDto]
	DeleteNoteTransaction(
		ctx context.Context,
		userBo userbos.UserBo,
		noteIdDto nDTOs.NoteIdDto,
	) single.Single[cDTOs.SuccessDto]
	GetNoteById(
		ctx context.Context,
		userBo userbos.UserBo,
		sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto],
	) single.Single[nDTOs.NoteReadDto]
	GetNotesPage(
		ctx context.Context,
		userBo userbos.UserBo,
		sessReqDto cDTOs.UKeySessionReqDto[pagination.PageRequest],
	) single.Single[pagination.Page[nDTOs.NotePreviewDto]]
	DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64]
}

type NoteServiceImpl struct {
	noteRepository repositories.NoteRepository
	userKeyService externalservices.ExtUserKeyService
	crudDSHandler  dshandlers.CrudDSHandler
	errorService   sharedservices.ErrorService
	noteBr         businessrules.NoteBr
}

func (n NoteServiceImpl) AddNoteTransaction(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto],
) single.Single[cDTOs.SuccessDto] {
	return dshandlers.TransactionalSingle(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) single.Single[cDTOs.SuccessDto] {
			sessDto, noteCreateDto := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
			keyDtoSrc := single.FromSupplierCached(func() (kDTOs.UserKeyDto, error) {
				return n.userKeyService.GetKeyFromSession(ctx, sessDto)
			})
			keySrc := single.MapWithError(keyDtoSrc, kDTOs.UserKeyDto.GetKey).ScheduleLazyAndCache(ctx)
			titleCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteCreateDto.Title))
			})
			textCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteCreateDto.Text))
			})
			noteSaveSrc := single.MapWithError(single.Zip3(keyDtoSrc, titleCipherSrc, textCipherSrc),
				func(t tuple.T3[kDTOs.UserKeyDto, []byte, []byte]) (models.Note, error) {
					keyDto, titleCipher, textCipher := t.V1, t.V2, t.V3
					note := models.Note{
						UserId:      userBo.Id,
						TextCipher:  textCipher,
						TitleCipher: titleCipher,
						KeyVersion:  keyDto.KeyVersion,
					}
					return n.noteRepository.Create(ctx, note)
				})
			return single.Map(noteSaveSrc, func(_ models.Note) cDTOs.SuccessDto {
				return cDTOs.NewSuccessTrue()
			})
		})

}

func (n NoteServiceImpl) UpdateNoteTransaction(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto],
) single.Single[cDTOs.SuccessDto] {
	return dshandlers.TransactionalSingle(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) single.Single[cDTOs.SuccessDto] {
			sessDto, noteUpdateDto := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
			existingSrc := n.getExistingNote(ctx, noteUpdateDto.Id).ScheduleLazyAndCache(ctx)
			keyDtoSrc := single.FromSupplierCached(func() (kDTOs.UserKeyDto, error) {
				return n.userKeyService.GetKeyFromSession(ctx, sessDto)
			})
			keySrc := single.MapWithError(single.Zip2(existingSrc, keyDtoSrc),
				func(t tuple.T2[models.Note, kDTOs.UserKeyDto]) ([]byte, error) {
					existing, keyDto := t.V1, t.V2
					err := n.noteBr.ValidateNoteUpdate(userBo, keyDto, existing)
					if err != nil {
						return nil, err
					}
					return keyDto.GetKey()
				}).ScheduleLazyAndCache(ctx)
			titleCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteUpdateDto.Title))
			})
			textCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteUpdateDto.Text))
			})
			noteSaveSrc := single.MapWithError(single.Zip4(existingSrc, keyDtoSrc, titleCipherSrc, textCipherSrc),
				func(t tuple.T4[models.Note, kDTOs.UserKeyDto, []byte, []byte]) (models.Note, error) {
					existing, keyDto, titleCipher, textCipher := t.V1, t.V2, t.V3, t.V4
					existing.TitleCipher = titleCipher
					existing.TextCipher = textCipher
					existing.KeyVersion = keyDto.KeyVersion
					return n.noteRepository.Update(ctx, existing)
				})
			return single.Map(noteSaveSrc, func(_ models.Note) cDTOs.SuccessDto {
				return cDTOs.NewSuccessTrue()
			})
		})
}

func (n NoteServiceImpl) DeleteNoteTransaction(
	ctx context.Context,
	userBo userbos.UserBo,
	noteIdDto nDTOs.NoteIdDto,
) single.Single[cDTOs.SuccessDto] {
	return dshandlers.TransactionalSingle(ctx, n.crudDSHandler,
		func(_ dshandlers.Session, ctx context.Context) single.Single[cDTOs.SuccessDto] {
			existingSrc := n.getExistingNote(ctx, noteIdDto.Id).ScheduleLazyAndCache(ctx)
			validateDeleteSrc := single.MapWithError(existingSrc,
				func(existing models.Note) (any, error) {
					err := n.noteBr.ValidateNoteDelete(userBo, existing)
					return any(true), err
				})
			noteDeleteSrc := single.MapWithError(single.Zip2(existingSrc, validateDeleteSrc),
				func(t tuple.T2[models.Note, any]) (models.Note, error) {
					existing := t.V1
					return n.noteRepository.Delete(ctx, existing)
				})
			return single.Map(noteDeleteSrc, func(_ models.Note) cDTOs.SuccessDto {
				return cDTOs.NewSuccessTrue()
			})
		})
}

func (n NoteServiceImpl) GetNoteById(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto],
) single.Single[nDTOs.NoteReadDto] {
	sessDto, noteIdDto := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
	existingSrc := n.getExistingNote(ctx, noteIdDto.Id).ScheduleLazyAndCache(ctx)
	keyDtoSrc := single.FromSupplierCached(func() (kDTOs.UserKeyDto, error) {
		return n.userKeyService.GetKeyFromSession(ctx, sessDto)
	})
	validationSrc := single.MapWithError(single.Zip2(existingSrc, keyDtoSrc),
		func(t tuple.T2[models.Note, kDTOs.UserKeyDto]) (any, error) {
			existing, keyDto := t.V1, t.V2
			return any(true), n.noteBr.ValidateNoteRead(userBo, keyDto, existing)
		})
	keySrc := single.MapWithError(single.Zip2(keyDtoSrc, validationSrc),
		func(t tuple.T2[kDTOs.UserKeyDto, any]) ([]byte, error) {
			keyDto := t.V1
			return keyDto.GetKey()
		},
	).ScheduleLazyAndCache(ctx)

	return single.FlatMap(single.Zip3(existingSrc, keySrc, validationSrc),
		func(t tuple.T3[models.Note, []byte, any]) single.Single[nDTOs.NoteReadDto] {
			existing, keyBytes := t.V1, t.V2
			textSrc := single.FromSupplierCached(func() (string, error) {
				txtBytes, err := cipherutils.DecryptAES(keyBytes, existing.TextCipher)
				return string(txtBytes), err
			})
			titleSrc := single.FromSupplierCached(func() (string, error) {
				titleBytes, err := cipherutils.DecryptAES(keyBytes, existing.TitleCipher)
				return string(titleBytes), err
			})
			return single.Map(single.Zip2(textSrc, titleSrc), func(t tuple.T2[string, string]) nDTOs.NoteReadDto {
				text, title := t.V1, t.V2
				coreNoteDetails := nDTOs.NewCoreNoteDetailsDto(title, text)
				noteDetailsDto := nDTOs.NoteReadDto{}
				mappers.MapCoreNoteDetailsAndNoteToNoteReadDto(&coreNoteDetails, &existing, &noteDetailsDto)
				return noteDetailsDto
			})
		})
}

func (n NoteServiceImpl) GetNotesPage(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[pagination.PageRequest],
) single.Single[pagination.Page[nDTOs.NotePreviewDto]] {
	sessionDto, pageRequest := sessReqDto.SetUserIdAndUnwrap(userBo.Id)
	err := n.noteBr.ValidateGetNotes(pageRequest)
	if err != nil {
		return single.Error[pagination.Page[nDTOs.NotePreviewDto]](err)
	}
	zippedSrc := single.FromSupplierCached(func() (tuple.T3[kDTOs.UserKeyDto, []byte, int64], error) {
		keyDtoSrc := single.FromSupplierCached(func() (kDTOs.UserKeyDto, error) {
			return n.userKeyService.GetKeyFromSession(ctx, sessionDto)
		})
		keySrc := single.MapWithError(keyDtoSrc, kDTOs.UserKeyDto.GetKey)
		countSrc := single.FromSupplierCached(func() (int64, error) {
			return n.noteRepository.CountByUserId(ctx, userBo.Id)
		})
		return single.RetrieveValue(ctx, single.Zip3(keyDtoSrc, keySrc, countSrc))
	})
	return single.FlatMap(zippedSrc,
		func(t tuple.T3[kDTOs.UserKeyDto, []byte, int64]) single.Single[pagination.Page[nDTOs.NotePreviewDto]] {
			keyDto, keyBytes, count := t.V1, t.V2, t.V3
			notes, err := n.noteRepository.GetPaginatedByUserId(ctx, userBo.Id, pageRequest)
			if err != nil {
				return single.Error[pagination.Page[nDTOs.NotePreviewDto]](err)
			}
			noteDTOs := make([]nDTOs.NotePreviewDto, 0, len(notes))
			for _, note := range notes {
				if note.KeyVersion != keyDto.KeyVersion {
					ruleErr := n.errorService.RuleErrorFromCode(apperrors.ErrCodeDataRace)
					return single.Error[pagination.Page[nDTOs.NotePreviewDto]](
						apperrors.NewBadReqErrorFromRuleError(ruleErr))
				}
				txtBytes, err := cipherutils.DecryptAES(keyBytes, note.TextCipher)
				if err != nil {
					return single.Error[pagination.Page[nDTOs.NotePreviewDto]](err)
				}
				titleBytes, err := cipherutils.DecryptAES(keyBytes, note.TitleCipher)
				if err != nil {
					return single.Error[pagination.Page[nDTOs.NotePreviewDto]](err)
				}
				txt, title := string(txtBytes), string(titleBytes)
				textPreview := utils.StringFirstNChars(txt, 60)
				coreNoteDto := nDTOs.NewCoreNoteDto(title)
				noteReadDto := nDTOs.NotePreviewDto{}
				mappers.MapTextPreviewAndCoreNoteAndNoteToNotePreviewDto(textPreview, &coreNoteDto, &note, &noteReadDto)
				noteDTOs = append(noteDTOs, noteReadDto)
			}
			noteDTOsSrc := single.Just(noteDTOs)
			return single.Map(noteDTOsSrc, func(noteDTOs []nDTOs.NotePreviewDto) pagination.Page[nDTOs.NotePreviewDto] {
				return pagination.NewPage(noteDTOs, count)
			})
		})
}

func (u NoteServiceImpl) DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64] {
	return single.FromSupplierCached(func() (int64, error) {
		return u.noteRepository.DeleteByUserIdAndGetCount(ctx, userId)
	})
}

func (n NoteServiceImpl) getExistingNote(ctx context.Context, id string) single.Single[models.Note] {
	findSrc := single.FromSupplierCached(func() (option.Maybe[models.Note], error) {
		return n.noteRepository.FindById(ctx, id)
	})
	return single.FlatMap(findSrc, func(m option.Maybe[models.Note]) single.Single[models.Note] {
		return option.Map(m, single.Just[models.Note]).
			OrElseGet(func() single.Single[models.Note] {
				ruleErr := n.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
				return single.Error[models.Note](apperrors.NewBadReqErrorFromRuleError(ruleErr))
			})
	})
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
