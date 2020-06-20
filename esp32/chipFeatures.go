package esp32

import "strings"

type Feature int
type Features map[Feature]bool
const (
    WiFi        Feature = iota
    Bluetooth
    SingleCore
    DualCore
    Clock160MHz
    Clock240MHz
    EmbeddedFlash
    VRefCalibrationEFuse
    BLK3Reserved
    CodingSchemeNone
    CodingScheme3_4
    CodingSchemeRepeat
    CodingSchemeInvalid
)

func (f Feature) String() string {
    return map[Feature]string{
        WiFi: "WiFi",
        Bluetooth: "Bluetooth",
        SingleCore: "Single Core",
        DualCore: "Dual Core",
        Clock160MHz: "160MHz",
        Clock240MHz: "240MHz",
        EmbeddedFlash: "Embedded Flash",
        VRefCalibrationEFuse: "VRef calibration in efuse",
        BLK3Reserved: "BLK3 partially reserved",
        CodingSchemeNone: "Coding Scheme None",
        CodingScheme3_4: "Coding Scheme 3/4",
        CodingSchemeRepeat: "Coding Scheme Repeat (UNSUPPORTED)",
        CodingSchemeInvalid: "Coding Scheme Invalid",
    }[f]
}

func (f Features) String() string {
    res := []string{}
    for feature,status := range f {
        if !status { continue }
        res = append(res, feature.String())
    }
    return strings.Join(res, ", ")
}

func (e *ESP32ROM) GetFeatures() (Features, error) {
    features := Features{
        WiFi: true,
    }

    word3, err := e.ReadEfuse(3)
    if err != nil { return features, err }

    features[Bluetooth] = word3[0] & (1 << 1) == 0
    features[DualCore] = word3[0] & (1 << 0) > 0
    features[SingleCore] = !features[DualCore]
    if word3[1] & (1 << 5) > 0 {
        features[Clock160MHz] = word3[1] & (1 << 4) > 0
        features[Clock240MHz] = !features[Clock160MHz]
    }

    pkgVersion := (word3[1] >> 1) & 0x07
    features[EmbeddedFlash] = pkgVersion == 2 || pkgVersion == 4 || pkgVersion == 5

    word4, err := e.ReadEfuse(4)
    if err != nil { return features, err }

    features[VRefCalibrationEFuse] = word4[1] & 0x1F > 0
    features[BLK3Reserved] = word4[1] >> 6 & 0x01 > 0

    word6, err := e.ReadEfuse(6)
    if err != nil { return features, err }

    features[CodingSchemeNone]    = word6[0] & 0x03 == 0
    features[CodingScheme3_4]     = word6[0] & 0x03 == 1
    features[CodingSchemeRepeat]  = word6[0] & 0x03 == 2
    features[CodingSchemeInvalid] = word6[0] & 0x03 == 3

    return features, nil
}
