package rom

import (
	"bytes"
	"fmt"
)

// collection modes
// i don't know what the difference between the two find modes is
const (
	CollectBuySatchel = 0x01
	CollectRingBox    = 0x02
	CollectUnderwater = 0x08
	CollectFind1      = 0x09
	CollectFind2      = 0x0a
	CollectAppear     = 0x1a // heart containers
	CollectFall       = 0x29
	CollectChest1     = 0x38 // most items
	CollectChest2     = 0x68 // map and compass
	CollectDig        = 0x5a
)

// A Treasure is data associated with a particular item ID and sub ID.
type Treasure struct {
	id, subID byte
	addr      uint16 // bank 15, value of hl at $15:466b

	// in order, starting at addr
	mode   byte // collection mode
	param  byte // parameter value to use for giveTreasure
	text   byte
	sprite byte
}

// SubID returns item sub ID of the treasure.
func (t Treasure) SubID() byte {
	return t.subID
}

func (t Treasure) CollectMode() byte {
	return t.mode
}

// RealAddr returns the total offset of the treasure data in a JP ROM.
func (t Treasure) RealAddr() int {
	return (&Addr{0x15, t.addr}).FullOffset()
}

// Bytes returns a slice of consecutive bytes of treasure data, as they would
// appear in the ROM.
func (t Treasure) Bytes() []byte {
	return []byte{t.mode, t.param, t.text, t.sprite}
}

// Mutate replaces the associated treasure in the given ROM data with this one.
func (t Treasure) Mutate(b []byte) error {
	// fake treasure
	if t.addr == 0 {
		return nil
	}

	addr, data := t.RealAddr(), t.Bytes()
	for i := 0; i < 4; i++ {
		b[addr+i] = data[i]
	}
	return nil
}

// Check verifies that the treasure's data matches the given ROM data.
func (t Treasure) Check(b []byte) error {
	addr, data := t.RealAddr(), t.Bytes()
	if bytes.Compare(b[addr:addr+4], data) != 0 {
		return fmt.Errorf("expected %x at %x; found %x",
			data, addr, b[addr:addr+4])
	}
	return nil
}

// Treasures maps item names to associated treasure data.
var Treasures = map[string]*Treasure{
	// equip items
	"shop shield L-1": &Treasure{0x01, 0x00, 0x52bd, 0x0a, 0x01, 0x1f, 0x13},
	"shield L-2":      &Treasure{0x01, 0x01, 0x52c1, 0x0a, 0x02, 0x20, 0x14},
	"bombs, 10":       &Treasure{0x03, 0x00, 0x52c9, 0x38, 0x10, 0x4d, 0x05},
	"sword 1":         &Treasure{0x05, 0x00, 0x52d9, 0x38, 0x01, 0x1c, 0x10},
	"sword 2":         &Treasure{0x05, 0x01, 0x52dd, 0x09, 0x01, 0x1c, 0x10},
	"boomerang L-1":   &Treasure{0x06, 0x00, 0x52f1, 0x0a, 0x01, 0x22, 0x1c},
	"boomerang L-2":   &Treasure{0x06, 0x01, 0x52f5, 0x38, 0x02, 0x23, 0x1d},
	"rod":             &Treasure{0x07, 0x00, 0x52f9, 0x38, 0x07, 0x0a, 0x1e},
	"spring":          &Treasure{0x07, 0x02, 0x5301, 0x09, 0x00, 0x0d, 0x1e},
	"summer":          &Treasure{0x07, 0x03, 0x5305, 0x09, 0x01, 0x0b, 0x1e},
	"autumn":          &Treasure{0x07, 0x04, 0x5309, 0x09, 0x02, 0x0c, 0x1e},
	"winter":          &Treasure{0x07, 0x05, 0x530d, 0x09, 0x03, 0x0a, 0x1e},
	"magnet gloves":   &Treasure{0x08, 0x00, 0x5149, 0x38, 0x00, 0x30, 0x18},
	"bombchus":        &Treasure{0x0d, 0x00, 0x531d, 0x0a, 0x10, 0x32, 0x24},
	"moosh's flute":   &Treasure{0x0e, 0x00, 0x5161, 0x0a, 0x0d, 0x3a, 0x4d},
	"dimitri's flute": &Treasure{0x0e, 0x00, 0x5161, 0x0a, 0x0c, 0x39, 0x4c},
	"strange flute":   &Treasure{0x0e, 0x00, 0x5161, 0x0a, 0x0d, 0x3b, 0x23},
	"ricky's flute":   &Treasure{0x0e, 0x00, 0x5161, 0x0a, 0x0b, 0x38, 0x4b},
	"slingshot 1":     &Treasure{0x13, 0x00, 0x5325, 0x38, 0x01, 0x2e, 0x21},
	"slingshot 2":     &Treasure{0x13, 0x01, 0x5329, 0x38, 0x01, 0x2e, 0x21},
	"shovel":          &Treasure{0x15, 0x00, 0x517d, 0x0a, 0x00, 0x25, 0x1b},
	"bracelet":        &Treasure{0x16, 0x00, 0x5181, 0x38, 0x00, 0x26, 0x19},
	"feather 1":       &Treasure{0x17, 0x00, 0x532d, 0x38, 0x01, 0x27, 0x16},
	"feather 2":       &Treasure{0x17, 0x01, 0x5331, 0x38, 0x01, 0x27, 0x16},
	"satchel 1":       &Treasure{0x19, 0x00, 0x52b5, 0x0a, 0x01, 0x2d, 0x20},
	"satchel 2":       &Treasure{0x19, 0x01, 0x52b9, 0x01, 0x01, 0x2d, 0x20},
	"fool's ore":      &Treasure{0x1e, 0x00, 0x51a1, 0x00, 0x00, 0x36, 0x4a},

	// not used because of progressive item upgrades
	// "sword L-2":       &Treasure{0x05, 0x01, 0x52dd, 0x09, 0x02, 0x1d, 0x11},
	// "slingshot L-2":   &Treasure{0x13, 0x01, 0x5329, 0x38, 0x02, 0x2f, 0x22},
	// "feather L-2":     &Treasure{0x17, 0x01, 0x5331, 0x38, 0x02, 0x28, 0x17},
	// "satchel 2":       &Treasure{0x19, 0x01, 0x52b9, 0x01, 0x00, 0x46, 0x20},

	// non-inventory items
	"rupees, 1":        &Treasure{0x28, 0x00, 0x5355, 0x38, 0x01, 0x01, 0x28},
	"rupees, 5":        &Treasure{0x28, 0x01, 0x5359, 0x38, 0x03, 0x02, 0x29},
	"rupees, 10":       &Treasure{0x28, 0x02, 0x535d, 0x38, 0x04, 0x03, 0x2a},
	"rupees, 20":       &Treasure{0x28, 0x03, 0x5361, 0x38, 0x05, 0x04, 0x2b},
	"rupees, 30":       &Treasure{0x28, 0x04, 0x5365, 0x38, 0x07, 0x05, 0x2b},
	"rupees, 50":       &Treasure{0x28, 0x05, 0x5369, 0x38, 0x0b, 0x06, 0x2c},
	"rupees, 100":      &Treasure{0x28, 0x06, 0x536d, 0x38, 0x0c, 0x07, 0x2d},
	"heart container":  &Treasure{0x2a, 0x00, 0x5399, 0x1a, 0x04, 0x16, 0x3b},
	"piece of heart":   &Treasure{0x2b, 0x01, 0x5391, 0x38, 0x01, 0x17, 0x3a},
	"rare peach stone": &Treasure{0x2b, 0x02, 0x5395, 0x02, 0x01, 0x17, 0x4e},

	// rings
	"discovery ring": &Treasure{0x2d, 0x04, 0x53c9, 0x38, 0x28, 0x54, 0x0e},
	"moblin ring":    &Treasure{0x2d, 0x05, 0x53cd, 0x38, 0x2b, 0x54, 0x0e},
	"steadfast ring": &Treasure{0x2d, 0x06, 0x53d1, 0x38, 0x10, 0x54, 0x0e},
	"rang ring L-1":  &Treasure{0x2d, 0x07, 0x53d5, 0x38, 0x0c, 0x54, 0x0e},
	"blast ring":     &Treasure{0x2d, 0x08, 0x53d9, 0x38, 0x0d, 0x54, 0x0e},
	"octo ring":      &Treasure{0x2d, 0x09, 0x53dd, 0x38, 0x2a, 0x54, 0x0e},
	"quicksand ring": &Treasure{0x2d, 0x0a, 0x53e1, 0x38, 0x23, 0x54, 0x0e},
	"armor ring L-2": &Treasure{0x2d, 0x0b, 0x53e5, 0x38, 0x05, 0x54, 0x0e},
	"power ring L-1": &Treasure{0x2d, 0x0e, 0x53f1, 0x38, 0x01, 0x54, 0x0e},
	"subrosian ring": &Treasure{0x2d, 0x10, 0x53f9, 0x38, 0x2d, 0x54, 0x0e},

	// dungeon items
	"small key":   &Treasure{0x30, 0x03, 0x5409, 0x38, 0x01, 0x1a, 0x42},
	"boss key":    &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"d1 boss key": &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"d2 boss key": &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"d3 boss key": &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"d6 boss key": &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"d7 boss key": &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"d8 boss key": &Treasure{0x31, 0x03, 0x5419, 0x38, 0x00, 0x1b, 0x43},
	"compass":     &Treasure{0x32, 0x02, 0x5425, 0x68, 0x00, 0x19, 0x41},
	"dungeon map": &Treasure{0x33, 0x02, 0x5431, 0x68, 0x00, 0x18, 0x40},

	// collection items
	"ring box L-1":    &Treasure{0x2c, 0x00, 0x53a5, 0x02, 0x01, 0x57, 0x33},
	"ring box L-2":    &Treasure{0x2c, 0x01, 0x53a9, 0x02, 0x02, 0x34, 0x34},
	"flippers":        &Treasure{0x2e, 0x00, 0x51e1, 0x02, 0x00, 0x31, 0x31},
	"gasha seed":      &Treasure{0x34, 0x01, 0x5341, 0x38, 0x01, 0x4b, 0x0d},
	"gnarled key":     &Treasure{0x42, 0x00, 0x5465, 0x29, 0x00, 0x42, 0x44},
	"floodgate key":   &Treasure{0x43, 0x00, 0x5235, 0x09, 0x00, 0x43, 0x45},
	"dragon key":      &Treasure{0x44, 0x00, 0x5239, 0x09, 0x00, 0x44, 0x46},
	"star ore":        &Treasure{0x45, 0x00, 0x523d, 0x5a, 0x00, 0x40, 0x57},
	"ribbon":          &Treasure{0x46, 0x00, 0x5241, 0x0a, 0x00, 0x41, 0x4f},
	"spring banana":   &Treasure{0x47, 0x00, 0x5245, 0x0a, 0x00, 0x66, 0x54},
	"ricky's gloves":  &Treasure{0x48, 0x00, 0x5249, 0x09, 0x01, 0x67, 0x55},
	"rusty bell":      &Treasure{0x4a, 0x00, 0x546d, 0x0a, 0x00, 0x55, 0x5b},
	"treasure map":    &Treasure{0x4b, 0x00, 0x5255, 0x0a, 0x00, 0x6c, 0x49},
	"round jewel":     &Treasure{0x4c, 0x00, 0x5259, 0x0a, 0x00, 0x47, 0x36},
	"pyramid jewel":   &Treasure{0x4d, 0x00, 0x5479, 0x08, 0x00, 0x4a, 0x37},
	"square jewel":    &Treasure{0x4e, 0x00, 0x5261, 0x38, 0x00, 0x48, 0x38},
	"x-shaped jewel":  &Treasure{0x4f, 0x00, 0x5265, 0x38, 0x00, 0x49, 0x39},
	"red ore":         &Treasure{0x50, 0x00, 0x5269, 0x38, 0x00, 0x3f, 0x59},
	"blue ore":        &Treasure{0x51, 0x00, 0x526d, 0x38, 0x00, 0x3e, 0x58},
	"hard ore":        &Treasure{0x52, 0x00, 0x5271, 0x0a, 0x00, 0x3d, 0x5a},
	"member's card":   &Treasure{0x53, 0x00, 0x5275, 0x0a, 0x00, 0x45, 0x48},
	"master's plaque": &Treasure{0x54, 0x00, 0x5279, 0x38, 0x00, 0x70, 0x26},

	// not real treasures, just placeholders for seeds in trees
	"ember tree seeds":   &Treasure{id: 0x00},
	"mystery tree seeds": &Treasure{id: 0x01},
	"scent tree seeds":   &Treasure{id: 0x02},
	"pegasus tree seeds": &Treasure{id: 0x03},
	"gale tree seeds 1":  &Treasure{id: 0x04},
	"gale tree seeds 2":  &Treasure{id: 0x05},
}

var seedIndexByTreeID = []byte{0, 4, 1, 2, 3, 3}

// FindTreasureName does a reverse lookup of the treasure in the map to return
// its name. It returns an empty string if not found.
func FindTreasureName(t *Treasure) string {
	for k, v := range Treasures {
		if v == t {
			return k
		}
	}
	return ""
}

// initialized automatically in init() based on contents of item slots
var TreasureIsUnique = map[string]bool{}

var uniqueIDTreasures = map[string]bool{}

func TreasureHasUniqueID(name string) bool {
	return uniqueIDTreasures[name]
}

// returns true iff a treasure can be lost permanently (i.e. outside of hide
// and seek).
func TreasureCanBeLost(name string) bool {
	switch name {
	case "shop shield L-1", "shield L-2", "star ore", "ribbon",
		"spring banana", "ricky's gloves", "round jewel", "pyramid jewel",
		"square jewel", "x-shapred jewel", "red ore", "blue ore", "hard ore":
		return true
	}
	return false
}
