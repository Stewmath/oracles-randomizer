package rom

import (
	"fmt"
	"log"
)

// A Mutable is a memory data that can be changed by the randomizer.
type Mutable interface {
	Mutate([]byte) error // change ROM bytes
	Check([]byte) error  // verify that the mutable matches the ROM
}

// A MutableRange is a length of mutable bytes starting at a given address.
type MutableRange struct {
	Addr     Addr
	Old, New []byte
}

// MutableByte returns a special case of MutableRange with a range of a single
// byte.
func MutableByte(addr Addr, old, new byte) *MutableRange {
	return &MutableRange{Addr: addr, Old: []byte{old}, New: []byte{new}}
}

// MutableWord returns a special case of MutableRange with a range of a two
// bytes.
func MutableWord(addr Addr, old, new uint16) *MutableRange {
	return &MutableRange{
		Addr: addr,
		Old:  []byte{byte(old >> 8), byte(old)},
		New:  []byte{byte(new >> 8), byte(new)},
	}
}

// Mutate replaces bytes in its range.
func (mr *MutableRange) Mutate(b []byte) error {
	addr := mr.Addr.FullOffset()
	for i, value := range mr.New {
		b[addr+i] = value
	}
	return nil
}

// Check verifies that the range matches the given ROM data.
func (mr *MutableRange) Check(b []byte) error {
	addr := mr.Addr.FullOffset()
	for i, value := range mr.Old {
		if b[addr+i] != value {
			return fmt.Errorf("expected %x at %x; found %x",
				mr.Old[i], addr+i, b[addr+i])
		}
	}
	return nil
}

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure            *Treasure
	IDAddrs, SubIDAddrs []Addr
	CollectMode         byte

	// TODO this is an incorrect model that happens to work for all currently
	//      slotted items except for the rod. for now the rod can have special
	//      logic, but this field really needs to be replaced with something
	//      more accurate (see treasureCollectionBehaviourTable in ages-disasm)
	SubIDOffset byte
}

// Mutate replaces the given IDs and subIDs in the given ROM data, and changes
// the associated treasure's collection mode as appropriate.
func (ms *MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.IDAddrs {
		b[addr.FullOffset()] = ms.Treasure.id
	}
	for _, addr := range ms.SubIDAddrs {
		// TODO see the comment on the SubIDOffset field of MutableSlot. for
		//      now, the rod needs special logic so it doesn't set an obtained
		//      season flag.
		if ms.SubIDOffset != 0 && ms.Treasure.id == 0x07 {
			b[addr.FullOffset()] = 0x07
		} else {
			b[addr.FullOffset()] = ms.Treasure.subID + ms.SubIDOffset
		}
	}
	ms.Treasure.mode = ms.CollectMode
	return ms.Treasure.Mutate(b)
}

// Check verifies that the slot's data matches the given ROM data.
func (ms *MutableSlot) Check(b []byte) error {
	for _, addr := range ms.IDAddrs {
		if b[addr.FullOffset()] != ms.Treasure.id {
			return fmt.Errorf("expected %x at %x; found %x",
				ms.Treasure.id, addr.FullOffset(), b[addr.FullOffset()])
		}
	}
	for _, addr := range ms.SubIDAddrs {
		if b[addr.FullOffset()] != ms.Treasure.subID+ms.SubIDOffset {
			return fmt.Errorf("expected %x at %x; found %x",
				ms.Treasure.subID+ms.SubIDOffset, addr.FullOffset(),
				b[addr.FullOffset()])
		}
	}
	if ms.CollectMode != ms.Treasure.mode {
		return fmt.Errorf("slot/treasure collect mode mismatch: %x/%x",
			ms.CollectMode, ms.Treasure.mode)
	}

	return nil
}

var ItemSlots = map[string]*MutableSlot{
	"d0 sword chest": &MutableSlot{
		Treasure:    Treasures["sword L-1"],
		IDAddrs:     []Addr{{0x0a, 0x7b86}},
		SubIDAddrs:  []Addr{{0x0a, 0x7b88}},
		SubIDOffset: 1,
		CollectMode: CollectChest,
	},
	"maku key fall": &MutableSlot{
		Treasure:    Treasures["gnarled key"],
		IDAddrs:     []Addr{{0x15, 0x657d}, {0x09, 0x7dff}, {0x09, 0x7de6}},
		SubIDAddrs:  []Addr{{0x15, 0x6580}, {0x09, 0x7e02}},
		CollectMode: CollectFall,
	},
	"boomerang gift": &MutableSlot{
		Treasure:    Treasures["boomerang L-1"],
		IDAddrs:     []Addr{{0x0b, 0x6648}},
		SubIDAddrs:  []Addr{{0x0b, 0x6649}},
		CollectMode: CollectFind2,
	},
	"rod gift": &MutableSlot{
		Treasure:    Treasures["rod"],
		IDAddrs:     []Addr{{0x15, 0x7511}},
		SubIDAddrs:  []Addr{{0x15, 0x750f}},
		SubIDOffset: 1,
		CollectMode: CollectChest, // it's what the data says
	},
	"shovel gift": &MutableSlot{
		Treasure:    Treasures["shovel"],
		IDAddrs:     []Addr{{0x0b, 0x6a6e}},
		SubIDAddrs:  []Addr{{0x0b, 0x6a6f}},
		CollectMode: CollectFind2,
	},
	"d1 satchel": &MutableSlot{
		// addresses are backwards from a normal slot
		Treasure:    Treasures["satchel"],
		IDAddrs:     []Addr{{0x09, 0x669b}},
		SubIDAddrs:  []Addr{{0x09, 0x669a}},
		CollectMode: CollectFind2,
	},
	"d2 bracelet chest": &MutableSlot{
		Treasure:    Treasures["bracelet"],
		IDAddrs:     []Addr{{0x15, 0x5424}},
		SubIDAddrs:  []Addr{{0x15, 0x5425}},
		CollectMode: CollectChest,
	},
	"blaino gift": &MutableSlot{
		Treasure:    Treasures["ricky's gloves"],
		IDAddrs:     []Addr{{0x0b, 0x64ce}},
		SubIDAddrs:  []Addr{{0x0b, 0x64cf}},
		CollectMode: CollectFind1,
	},
	"floodgate key gift": &MutableSlot{
		Treasure:    Treasures["floodgate key"],
		IDAddrs:     []Addr{{0x09, 0x626b}},
		SubIDAddrs:  []Addr{{0x09, 0x626a}},
		CollectMode: CollectFind1,
	},
	"square jewel chest": &MutableSlot{
		Treasure:    Treasures["square jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7397}},
		SubIDAddrs:  []Addr{{0x0b, 0x739b}},
		CollectMode: CollectChest,
	},
	"x-shaped jewel chest": &MutableSlot{
		Treasure:    Treasures["x-shaped jewel"],
		IDAddrs:     []Addr{{0x15, 0x53cd}},
		SubIDAddrs:  []Addr{{0x15, 0x53ce}},
		CollectMode: CollectChest,
	},
	"star ore spot": &MutableSlot{
		Treasure:    Treasures["star ore"],
		IDAddrs:     []Addr{{0x08, 0x62f4}, {0x08, 0x62fe}},
		SubIDAddrs:  []Addr{}, // special case, not set at all
		CollectMode: CollectDig,
	},
	"d3 feather chest": &MutableSlot{
		Treasure:    Treasures["feather L-1"],
		IDAddrs:     []Addr{{0x15, 0x5458}},
		SubIDAddrs:  []Addr{{0x15, 0x5459}},
		CollectMode: CollectChest,
	},
	"master's plaque chest": &MutableSlot{
		Treasure:    Treasures["master's plaque"],
		IDAddrs:     []Addr{{0x15, 0x554d}},
		SubIDAddrs:  []Addr{{0x15, 0x554e}},
		CollectMode: CollectChest,
	},
	"flippers gift": &MutableSlot{
		Treasure:    Treasures["flippers"],
		IDAddrs:     []Addr{{0x0b, 0x7310}, {0x0b, 0x72f3}},
		SubIDAddrs:  []Addr{{0x0b, 0x7311}},
		CollectMode: CollectFind2,
	},
	"spring banana tree": &MutableSlot{
		Treasure:    Treasures["spring banana"],
		IDAddrs:     []Addr{{0x09, 0x66b0}},
		SubIDAddrs:  []Addr{{0x09, 0x66af}},
		CollectMode: CollectFind2,
	},
	"dragon key spot": &MutableSlot{
		Treasure:    Treasures["dragon key"],
		IDAddrs:     []Addr{{0x09, 0x628d}},
		SubIDAddrs:  []Addr{{0x09, 0x628c}},
		CollectMode: CollectFind1,
	},
	"pyramid jewel spot": &MutableSlot{
		Treasure:    Treasures["pyramid jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7350}},
		SubIDAddrs:  []Addr{{0x0b, 0x7351}},
		CollectMode: CollectUnderwater,
	},
	// don't use this slot; no one knows about it and it's not required for
	// anything in a normal playthrough
	/*
		"ring box L-2 gift": &MutableSlot{
			Treasure:    Treasures["ring box L-2"],
			IDAddrs:     []Addr{{0x0b, 0x5c1a}},
			SubIDAddrs:  []Addr{{0x0b, 0x5c1b}},
			CollectMode: CollectGoronGift,
		},
	*/
	"d4 slingshot chest": &MutableSlot{
		Treasure:    Treasures["slingshot L-1"],
		IDAddrs:     []Addr{{0x15, 0x5470}},
		SubIDAddrs:  []Addr{{0x15, 0x5471}},
		CollectMode: CollectChest,
	},
	"d5 magnet gloves chest": &MutableSlot{
		Treasure:    Treasures["magnet gloves"],
		IDAddrs:     []Addr{{0x15, 0x5480}},
		SubIDAddrs:  []Addr{{0x15, 0x5481}},
		CollectMode: CollectChest,
	},
	"round jewel gift": &MutableSlot{
		Treasure:    Treasures["round jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7334}},
		SubIDAddrs:  []Addr{{0x0b, 0x7335}},
		CollectMode: CollectFind2,
	},
	"noble sword spot": &MutableSlot{
		// two cases depending on which sword you enter with
		Treasure:    Treasures["sword L-2"],
		IDAddrs:     []Addr{{0x0b, 0x6417}, {0x0b, 0x641e}},
		SubIDAddrs:  []Addr{{0x0b, 0x6418}, {0x0b, 0x641f}},
		CollectMode: CollectFind1,
	},
	"d6 boomerang chest": &MutableSlot{
		Treasure:    Treasures["boomerang L-2"],
		IDAddrs:     []Addr{{0x15, 0x54c0}},
		SubIDAddrs:  []Addr{{0x15, 0x54c1}},
		CollectMode: CollectChest,
	},
	"rusty bell spot": &MutableSlot{
		Treasure:    Treasures["rusty bell"],
		IDAddrs:     []Addr{{0x09, 0x6476}},
		SubIDAddrs:  []Addr{{0x09, 0x6475}},
		CollectMode: CollectFind2,
	},
	"d7 cape chest": &MutableSlot{
		Treasure:    Treasures["feather L-2"],
		IDAddrs:     []Addr{{0x15, 0x54e1}},
		SubIDAddrs:  []Addr{{0x15, 0x54e2}},
		CollectMode: CollectChest,
	},
	"d8 HSS chest": &MutableSlot{
		Treasure:    Treasures["slingshot L-2"],
		IDAddrs:     []Addr{{0x15, 0x551d}},
		SubIDAddrs:  []Addr{{0x15, 0x551e}},
		CollectMode: CollectChest,
	},

	// these are "fake" item slots in that they don't slot real treasures
	"ember tree": &MutableSlot{
		Treasure: Treasures["ember tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x64cb}},
	},
	"mystery tree": &MutableSlot{
		Treasure: Treasures["mystery tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x67dd}},
	},
	"scent tree": &MutableSlot{
		Treasure: Treasures["scent tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x685c}},
	},
	"pegasus tree": &MutableSlot{
		Treasure: Treasures["pegasus tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x6870}},
	},
	"sunken gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 1"],
		IDAddrs:  []Addr{{0x11, 0x69b0}},
	},
	"tarm gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 2"],
		IDAddrs:  []Addr{{0x11, 0x6a46}},
	},
}

// consider these mutables constants; they aren't changed in the randomization
// process.
var constMutables = map[string]Mutable{
	// have maku gate open from start
	"maku gate check": MutableByte(Addr{0x04, 0x61a3}, 0x7e, 0x66),

	// have horon village shop stock *and* sell items from the start, including
	// the flute. also don't disable the flute appearing until actually getting
	// ricky's flute; normally it disappears as soon as you enter the screen
	// northeast of d1 (or ricky's spot, whichever comes first).
	"horon shop stock check":   MutableByte(Addr{0x08, 0x4adb}, 0x05, 0x02),
	"horon shop sell check":    MutableByte(Addr{0x08, 0x48d0}, 0x05, 0x02),
	"horon shop flute check 1": MutableByte(Addr{0x08, 0x4b02}, 0xcb, 0xf6),
	"horon shop flute check 2": MutableByte(Addr{0x08, 0x4afc}, 0x6f, 0x7f),

	// subrosian dancing's flute prize is normally disabled by visiting the
	// same areas as the horon shop's flute.
	"dance hall flute check": MutableByte(Addr{0x09, 0x5e21}, 0x20, 0x80),

	// initiate all these events without requiring essences
	"ricky spawn check":         MutableByte(Addr{0x09, 0x4e68}, 0xcb, 0xf6),
	"dimitri essence check":     MutableByte(Addr{0x09, 0x4e36}, 0xcb, 0xf6),
	"dimitri flipper check":     MutableByte(Addr{0x09, 0x4e4c}, 0x2e, 0x04),
	"master essence check 1":    MutableByte(Addr{0x0a, 0x4bf5}, 0x02, 0x00),
	"master essence check 2":    MutableByte(Addr{0x0a, 0x4bea}, 0x40, 0x02),
	"master essence check 3":    MutableByte(Addr{0x08, 0x5887}, 0x40, 0x02),
	"round jewel essence check": MutableByte(Addr{0x0a, 0x4f8b}, 0x05, 0x00),
	"pirate essence check":      MutableByte(Addr{0x08, 0x6c32}, 0x20, 0x00),
	"eruption check 1":          MutableByte(Addr{0x08, 0x7c41}, 0x07, 0x00),
	"eruption check 2":          MutableByte(Addr{0x08, 0x7cd3}, 0x07, 0x00),

	// stop rosa from spawning and activate her portal by default. the first is
	// an essence check and the second is an edit to tile replacement data.
	"rosa spawn check": MutableByte(Addr{0x09, 0x678c}, 0x40, 0x04),
	"activate rosa portal": &MutableRange{Addr{0x04, 0x6016},
		[]byte{0x40, 0x33, 0xc5}, []byte{0x10, 0x33, 0xe6}},

	// count number of essences, not highest number essence
	"maku seed check 1": MutableByte(Addr{0x09, 0x7d8d}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7d8f}, 0x30, 0x18),

	// feather game: don't give fools ore, and don't return fools ore
	"get fools ore 1": MutableByte(Addr{0x14, 0x4111}, 0xe0, 0xf0),
	"get fools ore 2": MutableByte(Addr{0x14, 0x4112}, 0x2e, 0xf0),
	"get fools ore 3": MutableByte(Addr{0x14, 0x4113}, 0x5d, 0xf0),
	// There are tables indicating extra items to "get" and "lose" upon getting
	// an item. We remove the "lose fools ore" entry and insert a "get seeds
	// from slingshot" entry.
	"lose fools, get seeds from slingshot 1": MutableByte(Addr{0x3f, 0x4543}, 0x00, 0x13),
	"lose fools, get seeds from slingshot 2": &MutableRange{Addr{0x3f, 0x4545},
		[]byte{0x45, 0x00, 0x52, 0x50, 0x51, 0x17, 0x1e, 0x00},
		[]byte{0x20, 0x00, 0x46, 0x45, 0x00, 0x52, 0x50, 0x51}},
	"lose fools, get seeds from slingshot 3": MutableByte(Addr{0x3f, 0x44cf}, 0x44, 0x47),
	// since slingshot doesn't increment seed capacity, set the level-zero
	// capacity of seeds to 20, and move the pointer up by one byte.
	"satchel capacity": &MutableRange{Addr{0x3f, 0x4617},
		[]byte{0x20, 0x50, 0x99}, []byte{0x20, 0x20, 0x50}},
	"satchel capacity pointer": MutableByte(Addr{0x3f, 0x460e}, 0x16, 0x17),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7baf}, 0x05, 0x10),

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1": MutableWord(Addr{0x14, 0x4c67}, 0xe951, 0xa94f),
	"noble sword anim 2": MutableWord(Addr{0x14, 0x4e37}, 0x8364, 0xdf60),

	// getting the L-2 (or L-3) sword in the lost woods gives you two items;
	// one for the item itself and another that gives you the item and also
	// makes you do a spin slash animation. change the second ID bytes to a
	// fake item so that one slot doesn't give two items / the same item twice.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a}, 0x05, 0x10),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421}, 0x05, 0x10),

	// by default the cliff from sunken city to woods of winter is a one-way
	// door, which can lead to tricky softlock problems. until the routing
	// algorithm is capable of handling that kind of thing, the default season
	// for that area is just going to be spring, so that you can use the flower
	// to get back up.
	"cliff default season": MutableByte(Addr{0x01, 0x7e43}, 0x02, 0x00),

	// remove the snow piles in front of the shovel house so that shovel isn't
	// required not to softlock there (it's still required not to softlock in
	// hide and seek 2)
	"remove snow piles": MutableByte(Addr{0x24, 0x5dfe}, 0xd9, 0x04),
}

var mapIconByTreeID = []byte{0x15, 0x19, 0x16, 0x17, 0x18, 0x18}

// like the item slots, these are (usually) no-ops until the randomizer touches
// them.
var varMutables = map[string]Mutable{
	// map pop-up icons for seed trees
	"tarm gale tree map icon":   MutableByte(Addr{0x02, 0x6cb3}, 0x18, 0x18),
	"sunken gale tree map icon": MutableByte(Addr{0x02, 0x6cb6}, 0x18, 0x18),
	"scent tree map icon":       MutableByte(Addr{0x02, 0x6cb9}, 0x16, 0x16),
	"pegasus tree map icon":     MutableByte(Addr{0x02, 0x6cbc}, 0x17, 0x17),
	"mystery tree map icon":     MutableByte(Addr{0x02, 0x6cbf}, 0x19, 0x19),
	"ember tree map icon":       MutableByte(Addr{0x02, 0x6cc2}, 0x15, 0x15),

	// these scenes use specific item sprites not tied to treasure data
	"wooden sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x65f4},
		Old:  []byte{0x60, 0x00, 0x00},
		New:  []byte{0x60, 0x00, 0x00},
	},
	"rod graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6ba3},
		Old:  []byte{0x60, 0x10, 0x21},
		New:  []byte{0x60, 0x10, 0x21},
	},
	"noble sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6975},
		Old:  []byte{0x4e, 0x1a, 0x50},
		New:  []byte{0x4e, 0x1a, 0x50},
	},
	"master sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6978},
		Old:  []byte{0x4e, 0x1a, 0x40},
		New:  []byte{0x4e, 0x1a, 0x40},
	},

	// the satchel and slingshot should contain the type of seeds that grow on
	// the horon village tree.
	"satchel initial seeds":   MutableByte(Addr{0x3f, 0x453b}, 0x20, 0x20),
	"slingshot initial seeds": MutableByte(Addr{0x3f, 0x4544}, 0x46, 0x20),

	// the correct type of seed needs to be selected by default, otherwise the
	// player may be unable to use seeds when they only have one type. there
	// could also be serious problems with the submenu when they *do* obtain a
	// second type if the selection isn't either of them.
	//
	// this works by overwriting a couple of unimportant bytes in file
	// initialization.
	"satchel initial selection":   MutableWord(Addr{0x07, 0x418e}, 0xa210, 0xbe00),
	"slingshot initial selection": MutableWord(Addr{0x07, 0x419a}, 0x2e02, 0xbf00),

	// allow seed collection if you have a slingshot, by checking for the given
	// initial seed type
	"carry seeds in slingshot": MutableByte(Addr{0x10, 0x4b19}, 0x19, 0x20),
}

// get a collated map of all mutables
func getAllMutables() map[string]Mutable {
	slotMutables := make(map[string]Mutable)
	for k, v := range ItemSlots {
		slotMutables[k] = v
	}
	treasureMutables := make(map[string]Mutable)
	for k, v := range Treasures {
		treasureMutables[k] = v
	}

	mutableSets := []map[string]Mutable{
		constMutables,
		treasureMutables,
		slotMutables,
		varMutables,
	}

	// initialize master map w/ adequate capacity
	count := 0
	for _, set := range mutableSets {
		count += len(set)
	}
	allMutables := make(map[string]Mutable, count)

	// add mutables to master map
	for _, set := range mutableSets {
		for k, v := range set {
			if _, ok := allMutables[k]; ok {
				log.Fatalf("duplicate mutable key: %s", k)
			}
			allMutables[k] = v
		}
	}

	return allMutables
}
