package services

import (
	"context"
	"github.com/barweiss/go-tuple"
	"github.com/joamaki/goreactive/stream"
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
	) single.Single[nDTOs.NoteDetailsDto]
	GetNotesPage(
		ctx context.Context,
		userBo userbos.UserBo,
		sessReqDto cDTOs.UKeySessionReqDto[pagination.PageRequest],
	) single.Single[pagination.Page[nDTOs.NoteReadDto]]
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
			sessDto, noteCreateDto := sessReqDto.Session, sessReqDto.Value
			keyDtoSrc := n.userKeyService.GetKeyFromSession(ctx, sessDto).ScheduleLazyAndCache(ctx)
			keySrc := single.MapWithError(keyDtoSrc, kDTOs.UserKeyDto.GetKey).ScheduleLazyAndCache(ctx)
			titleCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteCreateDto.Title))
			})
			textCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteCreateDto.Text))
			})
			noteSaveSrc := single.FlatMap(single.Zip3(keyDtoSrc, titleCipherSrc, textCipherSrc),
				func(t tuple.T3[kDTOs.UserKeyDto, []byte, []byte]) single.Single[models.Note] {
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
			sessDto, noteUpdateDto := sessReqDto.Session, sessReqDto.Value
			existingSrc := n.getExistingNote(ctx, noteUpdateDto.Id).ScheduleLazyAndCache(ctx)
			keyDtoSrc := n.userKeyService.GetKeyFromSession(ctx, sessDto).ScheduleLazyAndCache(ctx)
			keySrc := single.FlatMap(single.Zip2(existingSrc, keyDtoSrc),
				func(t tuple.T2[models.Note, kDTOs.UserKeyDto]) single.Single[[]byte] {
					existing, keyDto := t.V1, t.V2
					noteUpdateValidationSrc := n.noteBr.ValidateNoteUpdate(userBo, keyDto, existing)
					return single.MapWithError(noteUpdateValidationSrc, func(_ any) ([]byte, error) {
						return keyDto.GetKey()
					})
				}).ScheduleLazyAndCache(ctx)
			titleCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteUpdateDto.Title))
			})
			textCipherSrc := single.MapWithError(keySrc, func(key []byte) ([]byte, error) {
				return cipherutils.EncryptAES(key, []byte(noteUpdateDto.Text))
			})
			noteSaveSrc := single.FlatMap(single.Zip4(existingSrc, keyDtoSrc, titleCipherSrc, textCipherSrc),
				func(t tuple.T4[models.Note, kDTOs.UserKeyDto, []byte, []byte]) single.Single[models.Note] {
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
			validateDeleteSrc := single.FlatMap(existingSrc,
				func(existing models.Note) single.Single[any] {
					return n.noteBr.ValidateNoteDelete(userBo, existing)
				})
			noteDeleteSrc := single.FlatMap(single.Zip2(existingSrc, validateDeleteSrc),
				func(t tuple.T2[models.Note, any]) single.Single[models.Note] {
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
) single.Single[nDTOs.NoteDetailsDto] {
	sessDto, noteIdDto := sessReqDto.Session, sessReqDto.Value
	existingSrc := n.getExistingNote(ctx, noteIdDto.Id).ScheduleLazyAndCache(ctx)
	keyDtoSrc := n.userKeyService.GetKeyFromSession(ctx, sessDto).ScheduleLazyAndCache(ctx)
	validationSrc := single.FlatMap(single.Zip2(existingSrc, keyDtoSrc),
		func(t tuple.T2[models.Note, kDTOs.UserKeyDto]) single.Single[any] {
			existing, keyDto := t.V1, t.V2
			return n.noteBr.ValidateNoteRead(userBo, keyDto, existing)
		})
	keySrc := single.MapWithError(single.Zip2(keyDtoSrc, validationSrc),
		func(t tuple.T2[kDTOs.UserKeyDto, any]) ([]byte, error) {
			keyDto := t.V1
			return keyDto.GetKey()
		},
	).ScheduleLazyAndCache(ctx)

	return single.FlatMap(single.Zip3(existingSrc, keySrc, validationSrc),
		func(t tuple.T3[models.Note, []byte, any]) single.Single[nDTOs.NoteDetailsDto] {
			existing, keyBytes := t.V1, t.V2
			textSrc := single.FromSupplierCached(func() (string, error) {
				txtBytes, err := cipherutils.DecryptAES(keyBytes, existing.TextCipher)
				return string(txtBytes), err
			})
			titleSrc := single.FromSupplierCached(func() (string, error) {
				titleBytes, err := cipherutils.DecryptAES(keyBytes, existing.TitleCipher)
				return string(titleBytes), err
			})
			return single.Map(single.Zip2(textSrc, titleSrc), func(t tuple.T2[string, string]) nDTOs.NoteDetailsDto {
				text, title := t.V1, t.V2
				coreNoteDetails := nDTOs.NewCoreNoteDetailsDto(title, text)
				noteDetailsDto := nDTOs.NoteDetailsDto{}
				mappers.MapCoreNoteDetailsAndNoteToNoteDetailsDto(&coreNoteDetails, &existing, &noteDetailsDto)
				return noteDetailsDto
			})
		})
}

func (n NoteServiceImpl) GetNotesPage(
	ctx context.Context,
	userBo userbos.UserBo,
	sessReqDto cDTOs.UKeySessionReqDto[pagination.PageRequest],
) single.Single[pagination.Page[nDTOs.NoteReadDto]] {
	sessionDto, pageRequest := sessReqDto.Session, sessReqDto.Value
	validationSrc := n.noteBr.ValidateGetNotes(pageRequest)
	zippedSrc := single.FlatMap(validationSrc, func(_ any) single.Single[tuple.T3[kDTOs.UserKeyDto, []byte, int64]] {
		keyDtoSrc := n.userKeyService.GetKeyFromSession(ctx, sessionDto)
		keySrc := single.MapWithError(keyDtoSrc, kDTOs.UserKeyDto.GetKey)
		countSrc := n.noteRepository.CountByUserId(ctx, userBo.Id)
		return single.Zip3(keyDtoSrc, keySrc, countSrc)
	})
	return single.FlatMap(zippedSrc,
		func(t tuple.T3[kDTOs.UserKeyDto, []byte, int64]) single.Single[pagination.Page[nDTOs.NoteReadDto]] {
			keyDto, keyBytes, count := t.V1, t.V2, t.V3
			findManyObs := n.noteRepository.GetPaginatedByUserId(ctx, userBo.Id, pageRequest)
			noteDetailsObs := stream.FlatMap(findManyObs, func(note models.Note) stream.Observable[nDTOs.NoteReadDto] {
				if note.KeyVersion != keyDto.KeyVersion {
					ruleErr := n.errorService.RuleErrorFromCode(apperrors.ErrCodeDataRace)
					return stream.Error(apperrors.NewBadReqErrorFromRuleError(ruleErr))
				}
				txtBytes, err := cipherutils.DecryptAES(keyBytes, note.TextCipher)
				if err != nil {
					return stream.Error(err)
				}
				titleBytes, err := cipherutils.DecryptAES(keyBytes, note.TitleCipher)
				if err != nil {
					return stream.Error(err)
				}
				txt, title := string(txtBytes), string(titleBytes)
				textPreview := utils.StringFirstNChars(txt, 60)
				coreNoteDto := nDTOs.NewCoreNoteDto(title)
				noteReadDto := nDTOs.NoteReadDto{}
				mappers.MapTextPreviewAndCoreNoteAndNoteToNoteReadDto(textPreview, &coreNoteDto, &note, &noteReadDto)
				return stream.Just(noteReadDto)
			})
			noteDTOsSrc := single.FromObservableAsList(noteDetailsObs)
			return single.Map(noteDTOsSrc, func(noteDTOs []nDTOs.NoteReadDto) pagination.Page[nDTOs.NoteReadDto] {
				return pagination.NewPage(noteDTOs, count)
			})
		})
}

func (n NoteServiceImpl) getExistingNote(ctx context.Context, id string) single.Single[models.Note] {
	findSrc := n.noteRepository.FindById(ctx, id)
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
