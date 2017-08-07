package srec

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	TypeFieldStrLen   = 2
	LengthFieldStrLen = 2
	CSumFieldStrLen   = 2
)

type Srec struct {
	HeaderRecord  HeaderRecord
	BinaryRecords []BinaryRecord
	FooterRecord  FooterRecord
	StartAddress  uint
	EndAddress    uint
	Bytes         []byte
	OutStream     io.Writer
	ErrStream     io.Writer
}

type HeaderRecord struct {
	Length   uint32
	Data     []byte
	Checksum byte
}

type BinaryRecord struct {
	Srectype string
	Length   uint32
	Address  uint32
	Data     []byte
	Checksum byte
}

type FooterRecord struct {
	Srectype  string
	EntryAddr uint32
	Checksum  byte
}

var ()

func NewSrec(outs, errs io.Writer) *Srec {
	return &Srec{OutStream: outs, ErrStream: errs}
}

func NewHeaderRecord() *HeaderRecord {
	return &HeaderRecord{}
}

func NewBianryRecord() *BinaryRecord {
	return &BinaryRecord{}
}

func NewFooterRecord() *FooterRecord {
	return &FooterRecord{}
}

func (srs *Srec) ParseFile(fileReader io.Reader) {
	rec := new(BinaryRecord)

	scanner := bufio.NewScanner(fileReader)

	for scanner.Scan() {
		line := scanner.Text()
		sl := strings.Split(line, "")

		srectype := strings.Join(sl[:2], "")
		switch {
		case srectype == "S0":
		case (srectype == "S1") || (srectype == "S2") || (srectype == "S3"):
			rec.getSrecBinaryRecordFields(srectype, sl)
			srs.BinaryRecords = append(srs.BinaryRecords, *rec)
		case (srectype == "S7") || (srectype == "S8") || (srectype == "S9"):
		default:
			// pass S4~6
		}
	}
}

func (rec *BinaryRecord) getSrecBinaryRecordFields(srectype string, sl []string) error {
	var len uint64
	var addr uint64
	var data []byte
	var csum uint64

	addrStrLen, err := getAddrStrLen(srectype)
	if err != nil {
		return err
	}
	lengthStrLen, err := getLengthStrLen(sl)

	len, _ = strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	addr, _ = strconv.ParseUint(strings.Join(sl[4:4+addrStrLen], ""), 16, 32)
	data = make([]byte, 0)
	for i := (TypeFieldStrLen + LengthFieldStrLen + addrStrLen); i < (TypeFieldStrLen+LengthFieldStrLen)+(lengthStrLen-CSumFieldStrLen); i += 2 {
		b, _ := strconv.ParseUint(strings.Join(sl[i:i+2], ""), 16, 32)
		data = append(data, byte(b))
	}
	csum, _ = strconv.ParseUint(strings.Join(sl[4+(int(len)*2)-2:(4+(int(len)*2)-2)+2], ""), 16, 32)

	rec.Srectype = srectype
	rec.Length = uint32(len)
	rec.Address = uint32(addr)
	rec.Data = data
	rec.Checksum = byte(csum)
	return nil
}

func getAddrStrLen(srectype string) (int, error) {
	switch srectype {
	case "S1":
		return 4, nil
	case "S2":
		return 6, nil
	case "S3":
		return 8, nil
	default:
		return 0, fmt.Errorf("%s is not srectype", srectype)
	}
}

func getLengthStrLen(sl []string) (int, error) {
	len, err := strconv.ParseUint(strings.Join(sl[2:4], ""), 16, 32)
	return int(len * 2), err
}

func (sr *Srec) GetBytes(ByteSize uint32) []byte {
	var bytes []byte
	for _, br := range sr.BinaryRecords {
		for _, b := range br.Data {
			bytes = append(bytes, b)
			if uint32(len(bytes)) == ByteSize {
				return bytes
			}
		}
	}
	return bytes
}
