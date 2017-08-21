package srec

import "fmt"

func (sr *Srec) Bytes() []byte {
	return sr.dataBytes
}

func (sr *Srec) SetBytes(wAddr uint32, wBytes []byte) error {
	if len(sr.dataRecords) == 0 {
		return fmt.Errorf("byte data is empty. call PaeseFile() or maybe srec file has no S1~3 records.")
	}
	if (wAddr < sr.startAddress) || (wAddr > sr.endAddress) {
		return fmt.Errorf("data address 0x%08X is out of srec range.", wAddr)
	}
	start := int(wAddr) - int(sr.startAddress)
	for i := 0; i < len(wBytes); i++ {
		sr.dataBytes[start+i] = wBytes[i]
	}
	return nil
}
