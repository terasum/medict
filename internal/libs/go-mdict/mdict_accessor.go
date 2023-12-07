package go_mdict

import "encoding/json"

type MdictAccessor struct {
	Filepath          string `json:"filepath"`
	IsRecordEncrypted bool   `json:"is_record_encrypted"`
	IsMDD             bool   `json:"is_mdd"`
	IsUTF16           bool   `json:"is_utf_16"`
}

func NewAccessor(mdict *Mdict) *MdictAccessor {
	return &MdictAccessor{
		Filepath:          mdict.filePath,
		IsRecordEncrypted: mdict.meta.encryptType == EncryptRecordEnc,
		IsMDD:             mdict.fileType == MdictTypeMdd,
		IsUTF16:           mdict.meta.encoding == EncodingUtf16,
	}
}

func NewAccessorFromJson(data []byte) (*MdictAccessor, error) {
	mdi := new(MdictAccessor)
	err := json.Unmarshal(data, mdi)
	return mdi, err
}

func (mdi *MdictAccessor) Serialize() ([]byte, error) {
	return json.Marshal(mdi)
}

func (mdi *MdictAccessor) RetrieveDefByIndex(index *MDictKeywordIndex) ([]byte, error) {
	return locateDefByKWIndex(index, mdi.Filepath, mdi.IsRecordEncrypted, mdi.IsMDD, mdi.IsUTF16)
}
