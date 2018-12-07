package goinsight

import (
	"fmt"
	"unsafe"
)

type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type rtype struct {
	size       uintptr
	ptrdata    uintptr  // number of bytes in the type that can contain pointers
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      uint8    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        int32    // string form
	ptrToThis  int32    // type for pointer to this type, may be zero
}

type emptyInterface struct {
	typ  *rtype			// pointer to data type descripter
	word unsafe.Pointer // pointer to data
}

type iTab struct {
	ityp *rtype // static interface type, pointer to data type descripter of interface
	typ  *rtype // dynamic concrete type, pointer to data type descripter of data
	hash uint32 // copy of typ.hash
	_    [4]byte
	fun  [1]uintptr // method table
}

type nonEmptyInterface struct {
	// see ../runtime/iface.go:/Itab
	itab *iTab			// pointer to itab
	word unsafe.Pointer	// pointer to data
}

// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [8]uint8
	// Followed by bucketCnt keys and then bucketCnt values.
	// NOTE: packing all the keys together and then all the values together makes the
	// code a bit more complicated than alternating key/value/key/value/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

// mapextra holds fields that are not present on all maps.
type mapextra struct {
	// If both key and value do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.overflow and h.map.oldoverflow.
	// overflow and oldoverflow are only used if key and value do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}

type stringMapCts struct{
	hash uint64
	Kkeys [8]string
	Values [8]string
	plink unsafe.Pointer
}

type hmap struct {
	// Note: the format of the Hmap is encoded in ../../cmd/internal/gc/reflect.go and
	// ../reflect/type.go. Don't change this structure without also changing that code!
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    *[2]stringMapCts // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}

type rmap struct {
	 mm *hmap
}

// empty interface
func InsightEmptyInterface(i interface{}) {
	ei := *(*emptyInterface)(unsafe.Pointer(&i))
	fmt.Printf("**** Insight an empty interface from address : %p ****\n",&ei)
	fmt.Printf("Level1(addr=%p,size=%d) : %#v\n", &ei, unsafe.Sizeof(ei), ei)
	fmt.Printf("Level2(addr=%p,size=%d) : typ=%#v\n", ei.typ, unsafe.Sizeof(*(ei.typ)), *(ei.typ))
	fmt.Printf("Level3(addr=%p,size=%d) : typ.alg=%#v\n", ei.typ.alg, unsafe.Sizeof(*(ei.typ.alg)), *(ei.typ.alg))

}

// non empty interface
func InsightNonEmptyInterface(p unsafe.Pointer) {
	ei := *(*nonEmptyInterface)(p)
	fmt.Printf("**** Insight an non-empty interface from address : %p ****\n",p)
	fmt.Printf("Level1(addr=%p,size=%d) : %#v\n", &ei, unsafe.Sizeof(ei), ei)
	fmt.Printf("Level2(addr=%p,size=%d) : itab=%#v\n", ei.itab, unsafe.Sizeof(*(ei.itab)), *(ei.itab))
	fmt.Printf("Level3(addr=%p,size=%d) : itab.ityp=%#v\n", ei.itab.ityp, unsafe.Sizeof(*(ei.itab.ityp)), *(ei.itab.ityp))
	fmt.Printf("Level3(addr=%p,size=%d) : itab.typ=%#v\n", ei.itab.typ, unsafe.Sizeof(*(ei.itab.typ)), *(ei.itab.typ))
}

// map[string]string
func InsightMapString(a map[string]string){

	em := *(*rmap)(unsafe.Pointer(&a))
	fmt.Printf("**** Insight an map[string]string from address : %p ****\n",&em)
	fmt.Printf("Level1(addr=%p,size=%d) : %#v\n", &em, unsafe.Sizeof(em), em)
	fmt.Printf("Level2(addr=%p,size=%d) : %#v\n", em.mm, unsafe.Sizeof(*(em.mm)), *(em.mm))
	fmt.Printf("Level3(addr=%p,size=%d) : %#v\n", em.mm.buckets,unsafe.Sizeof(*(em.mm.buckets)),*(em.mm.buckets))
}

// insight memory values from the head address of an object
// method : p64,i64,i32,i16,i8,u64,u32,u16,u8,c8
// return : when method is p64, then add the pointer value to return
func InsightMem(i interface{}, method ...string) []unsafe.Pointer {
	ei := *(*emptyInterface)(unsafe.Pointer(&i))
	pt := ei.word
	ptnum := uintptr(pt)
	var rts []unsafe.Pointer
	for _, v := range method {
		switch v {
		case "p64":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[8] : %p\n", pt, *((*(unsafe.Pointer))(pt)))
			rts = append(rts, *((*(unsafe.Pointer))(pt)))
			ptnum += 8
		case "i8":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[1] : %d\n", pt, *((*int8)(pt)))
			ptnum += 1
		case "i16":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[2] : %d\n", pt, *((*int16)(pt)))
			ptnum += 2
		case "i32":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[4] : %d\n", pt, *((*int32)(pt)))
			ptnum += 4
		case "i64":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[8] : %d\n", pt, *((*int64)(pt)))
			ptnum += 8
		case "u8":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[1] : %d\n", pt, *((*uint8)(pt)))
			ptnum += 1
		case "u16":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[2] : %d\n", pt, *((*uint16)(pt)))
			ptnum += 2
		case "u32":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[4] : %d\n", pt, *((*uint32)(pt)))
			ptnum += 4
		case "u64":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[8] : %d\n", pt, *((*uint64)(pt)))
			ptnum += 8
		case "c8":
			pt = unsafe.Pointer(ptnum)
			fmt.Printf("%v[1] : %c\n", pt, *((*byte)(pt)))
			ptnum += 1
		}
	}
	return rts
}

func InsightMemString(a *string){
	InsightMem(*a,"p64","i64")
}

func InsightMemArray(a *[4]string){
	InsightMem(*a,"p64","i64","p64","i64","p64","i64","p64","i64")
}

func InsightMemSlice(i interface{}){
	ei := *(*emptyInterface)(unsafe.Pointer(&i))
	InsightMem(ei.word,"p64","i64","i64")
}

