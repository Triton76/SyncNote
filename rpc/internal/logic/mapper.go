package logic

import (
	"SyncNote/model"
	"SyncNote/rpc/pb/syncnoterpc"
)

func toNoteResp(note *model.Note) *syncnoterpc.NoteResp {
	if note == nil {
		return nil
	}
	return &syncnoterpc.NoteResp{
		NoteId:       note.NoteID,
		UserId:       note.UserID,
		Title:        note.Title,
		Content:      note.Content,
		Version:      note.Version,
		LastModified: note.LastModified.Unix(),
	}
}

func toNoteSummary(note *model.Note) *syncnoterpc.NoteSummary {
	if note == nil {
		return nil
	}
	return &syncnoterpc.NoteSummary{
		NoteId:       note.NoteID,
		Title:        note.Title,
		Version:      note.Version,
		LastModified: note.LastModified.Unix(),
	}
}

func toPbSaveCode(code model.SaveCode) syncnoterpc.SaveCode {
	switch code {
	case model.SaveCodeOK:
		return syncnoterpc.SaveCode_SAVE_CODE_OK
	case model.SaveCodeVersionConflict:
		return syncnoterpc.SaveCode_SAVE_CODE_VERSION_CONFLICT
	case model.SaveCodeNotFound:
		return syncnoterpc.SaveCode_SAVE_CODE_NOT_FOUND
	case model.SaveCodeInvalidParam:
		return syncnoterpc.SaveCode_SAVE_CODE_INVALID_PARAM
	default:
		return syncnoterpc.SaveCode_SAVE_CODE_UNSPECIFIED
	}
}

func toSaveNoteResp(result *model.SaveResult) *syncnoterpc.SaveNoteResp {
	if result == nil {
		return &syncnoterpc.SaveNoteResp{
			Success: false,
			Code:    syncnoterpc.SaveCode_SAVE_CODE_UNSPECIFIED,
			Message: "empty save result",
		}
	}

	resp := &syncnoterpc.SaveNoteResp{
		Success:       result.Success,
		Code:          toPbSaveCode(result.Code),
		Message:       result.Message,
		LatestVersion: result.LatestVersion,
		LatestContent: result.LatestContent,
	}
	if result.Note != nil {
		resp.Note = toNoteResp(result.Note)
	}

	return resp
}

func invalidParamSaveResp(message string) *syncnoterpc.SaveNoteResp {
	return &syncnoterpc.SaveNoteResp{
		Success: false,
		Code:    syncnoterpc.SaveCode_SAVE_CODE_INVALID_PARAM,
		Message: message,
	}
}
