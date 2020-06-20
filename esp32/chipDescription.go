package esp32

import "fmt"

type ChipDescription struct {
    ChipType    ChipType
    Revision    byte
}

func (c *ChipDescription) String() string {
    return fmt.Sprintf("%s (revision %d)", c.ChipType.String(), c.Revision)
}

func (e *ESP32ROM) GetChipDescription() (*ChipDescription, error) {
    word3, err := e.ReadEfuse(3)
    if err != nil {
        return nil, err
    }
    word5, err := e.ReadEfuse(5)
    if err != nil {
        return nil, err
    }
    apbCtlBase, err := e.ReadRegister(drRegSysconBase + 0x7C)

    revisionBit0 := (word3[1] >> 7) & 0x01
    revisionBit1 := (word5[2] >> 4) & 0x01
    revisionBit2 := (apbCtlBase[3] >> 7) & 0x01

    revision := byte(0)
    if revisionBit0 > 0 {
        if revisionBit1 > 0 {
            if revisionBit2 > 0 {
                revision = 3
            } else {
                revision = 2
            }
        } else {
            revision = 1
        }
    }

    return &ChipDescription{
        ChipType: ChipType((word3[1] >> 1) & 0x07),
        Revision: revision,
    }, nil
}
