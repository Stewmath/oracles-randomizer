package prenode

// This package contains definitions of nodes and node relationships before
// they are inserted into the graph. This is necessary because nodes
// relationships can't be made until the nodes are added first (and it's nice
// not to clutter the other packages with all these definitions).

// XXX need to be careful about rings. i can't imagine a situation where you'd
//     need both energy ring and fist ring, but if you did, then you'd need to
//     have the L-2 ring box to do so without danger of soft locking.

// A Type identifies whether a prenode is and And, Or, or Root node, whether it
// is an item slot, and whether it is a non-item slot milestone.
type Type int

// And, Or, and Root are pretty self-explanatory. One with a Slot suffix is an
// item slot, and one with a Step suffix is treated as a milestone for routing
// purposes. Slot types are also treated as steps; see the Point.IsStep()
// function.
//
// The following function are half syntactic sugar for declaring large lists of
// node relationships.
const (
	RootType Type = iota
	AndType
	OrType
	AndSlotType
	OrSlotType
	AndStepType
	OrStepType
)

// A Prenode is a mapping of strings that will become And or Or nodes in the
// graph. A prenode can have nested prenodes as parents instead of strings.
type Prenode struct {
	Parents []interface{}
	Type    Type
}

// CreateFunc returns a function that creates graph nodes from a list of key
// strings or sub-prenodes, based on the given prenode type.
func CreateFunc(prenodeType Type) func(parents ...interface{}) *Prenode {
	return func(parents ...interface{}) *Prenode {
		return &Prenode{Parents: parents, Type: prenodeType}
	}
}

// Convenience functions for creating prenodes succinctly. See the Type const
// comment for information on the various types.
var (
	Root    = CreateFunc(RootType)
	And     = CreateFunc(AndType)
	AndSlot = CreateFunc(AndSlotType)
	AndStep = CreateFunc(AndStepType)
	Or      = CreateFunc(OrType)
	OrSlot  = CreateFunc(OrSlotType)
	OrStep  = CreateFunc(OrStepType)
)

// BaseItems returns a map of item prenodes that may be assigned to slots.
func BaseItems() map[string]*Prenode {
	return baseItemPrenodes
}

// GetAll returns all prenodes.
func GetAll() map[string]*Prenode {
	total := make(map[string]*Prenode)
	appendPrenodes(total,
		itemPrenodes, baseItemPrenodes, ignoredBaseItemPrenodes, killPrenodes,
		holodrumPrenodes, subrosiaPrenodes, portalPrenodes,
		d0Prenodes, d1Prenodes, d2Prenodes, d3Prenodes, d4Prenodes,
		d5Prenodes, d6Prenodes, d7Prenodes, d8Prenodes, d9Prenodes)
	return total
}

// merge the given maps into the first argument
func appendPrenodes(total map[string]*Prenode, maps ...map[string]*Prenode) {
	for _, prenodeMap := range maps {
		for k, v := range prenodeMap {
			total[k] = v
		}
	}
}
