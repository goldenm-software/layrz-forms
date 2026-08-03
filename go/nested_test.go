package layrz

import (
	"testing"
)

// TestNestedSubformMinLength tests a non-nil subform with validation errors.
func TestNestedSubformMinLength(t *testing.T) {
	type Address struct {
		StreetName *string `layrz:"char,required,min_length=5"`
	}

	type Form struct {
		Address *Address `layrz:"subform"`
	}

	streetName := "abc" // 3 chars, less than min_length=5
	form := &Form{
		Address: &Address{
			StreetName: &streetName,
		},
	}

	errs := Validate(form)

	// Should have address.streetName error with minLength code
	addrErrs, ok := errs["address.streetName"]
	if !ok {
		t.Fatalf("expected address.streetName key, got keys: %v", errs.Keys())
	}
	if len(addrErrs) != 1 {
		t.Errorf("expected 1 error, got %d", len(addrErrs))
	}
	if addrErrs[0].Code != testCodeMinLength {
		t.Errorf("expected minLength, got %q", addrErrs[0].Code)
	}
	if addrErrs[0].Expected != 5 {
		t.Errorf("expected Expected=5, got %v", addrErrs[0].Expected)
	}
	if addrErrs[0].Received != 3 {
		t.Errorf("expected Received=3, got %v", addrErrs[0].Received)
	}
}

// TestNestedSubformListErrors tests a subform list with multiple validation errors.
func TestNestedSubformListErrors(t *testing.T) {
	type Item struct {
		Name  *string  `layrz:"char,required"`
		Price *float64 `layrz:"number,required,min_value=0"`
	}

	type Form struct {
		Items []*Item `layrz:"subform_list"`
	}

	form := &Form{
		Items: []*Item{
			{
				Name:  nil, // Missing required
				Price: Ptr(10.0),
			},
			{
				Name:  Ptr("Item 2"),
				Price: Ptr(-5.0), // Negative, violates min_value
			},
		},
	}

	errs := Validate(form)

	// Should have items.0.name error
	if _, ok := errs["items.0.name"]; !ok {
		t.Fatalf("expected items.0.name key, got keys: %v", errs.Keys())
	}

	// Should have items.1.price error
	if _, ok := errs["items.1.price"]; !ok {
		t.Fatalf("expected items.1.price key, got keys: %v", errs.Keys())
	}

	// Check the price error details
	priceErrs := errs["items.1.price"]
	if len(priceErrs) != 1 || priceErrs[0].Code != "minValue" {
		if len(priceErrs) == 0 {
			t.Errorf("expected minValue error for items.1.price, got 0 errors")
		} else {
			t.Errorf("expected minValue error for items.1.price, got code=%q expected=%v received=%v", priceErrs[0].Code, priceErrs[0].Expected, priceErrs[0].Received)
		}
	}
}

// TestNilSubformIsSkipped tests that a nil subform produces no errors.
func TestNilSubformIsSkipped(t *testing.T) {
	type Address struct {
		StreetName *string `layrz:"char,required"`
	}

	type Form struct {
		Address *Address `layrz:"subform"`
	}

	form := &Form{
		Address: nil, // Nil subform should be skipped entirely
	}

	errs := Validate(form)

	// Should have no errors
	if len(errs) > 0 {
		t.Errorf("expected no errors for nil subform, got %v", errs)
	}

	// Specifically, no address.* keys
	for key := range errs {
		if len(key) > 8 && key[:8] == "address." {
			t.Errorf("unexpected error for nil subform: %q", key)
		}
	}
}

// TestNestedSubformWithOwnCleanMethods tests that nested subforms' clean methods run with prefixed keys.
type NestedAddressWithClean struct {
	StreetName *string `layrz:"char,required"`
}

func (a *NestedAddressWithClean) CleanStreetClean() Errors {
	return Errors{
		"street_clean_check": {
			{Code: "customAddressError"},
		},
	}
}

type FormWithNestedClean struct {
	Address *NestedAddressWithClean `layrz:"subform"`
}

func TestNestedSubformWithOwnCleanMethods(t *testing.T) {
	form := &FormWithNestedClean{
		Address: &NestedAddressWithClean{
			StreetName: Ptr("valid street"),
		},
	}

	errs := Validate(form)

	// Should have address.streetCleanCheck error (from nested clean method, prefixed)
	if _, ok := errs["address.streetCleanCheck"]; !ok {
		t.Fatalf("expected address.streetCleanCheck key, got keys: %v", errs.Keys())
	}
}

// TestDeeplyNestedThreeLevels tests 3-level nesting with prefix composition.
type Level3 struct {
	Value *string `layrz:"char,required"`
}

type Level2Deep struct {
	L3 *Level3 `layrz:"subform"`
}

type Level1Deep struct {
	L2 *Level2Deep `layrz:"subform"`
}

type Level0Deep struct {
	L1 *Level1Deep `layrz:"subform"`
}

func TestDeeplyNestedThreeLevels(t *testing.T) {
	form := &Level0Deep{
		L1: &Level1Deep{
			L2: &Level2Deep{
				L3: &Level3{
					Value: nil, // Missing required
				},
			},
		},
	}

	errs := Validate(form)

	// Should have l1.l2.l3.value key
	if _, ok := errs["l1.l2.l3.value"]; !ok {
		t.Fatalf("expected l1.l2.l3.value key, got keys: %v", errs.Keys())
	}
}

// TestMultipleItemsInList tests a list with multiple items at different indices.
func TestMultipleItemsInList(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type Form struct {
		Items []*Item `layrz:"subform_list"`
	}

	form := &Form{
		Items: []*Item{
			{Name: Ptr("item0")},
			{Name: nil}, // Missing
			{Name: Ptr("item2")},
			{Name: nil}, // Missing
			{Name: Ptr("item4")},
		},
	}

	errs := Validate(form)

	// Should have items.1.name and items.3.name errors
	if _, ok := errs["items.1.name"]; !ok {
		t.Fatal("missing items.1.name")
	}
	if _, ok := errs["items.3.name"]; !ok {
		t.Fatal("missing items.3.name")
	}

	// Should NOT have items.0.name, items.2.name, items.4.name
	if _, ok := errs["items.0.name"]; ok {
		t.Error("unexpected items.0.name error")
	}
	if _, ok := errs["items.2.name"]; ok {
		t.Error("unexpected items.2.name error")
	}
	if _, ok := errs["items.4.name"]; ok {
		t.Error("unexpected items.4.name error")
	}
}

// TestNestedListOfNonPointerStructs tests subform_list with non-pointer struct elements.
func TestNestedListOfNonPointerStructs(t *testing.T) {
	type Item struct {
		Name *string `layrz:"char,required"`
	}

	type Form struct {
		Items []Item `layrz:"subform_list"` // Note: not []*Item
	}

	form := &Form{
		Items: []Item{
			{Name: Ptr("item0")},
			{Name: nil}, // Missing required
		},
	}

	errs := Validate(form)

	// Should have items.1.name error
	if _, ok := errs["items.1.name"]; !ok {
		t.Fatalf("expected items.1.name key, got keys: %v", errs.Keys())
	}
}

// TestNestedSubformCamelCasing tests that nested subform keys use camelCase.
func TestNestedSubformCamelCasing(t *testing.T) {
	type UserInfo struct {
		FirstName *string `layrz:"char,required"`
	}

	type PersonForm struct {
		UserInfo *UserInfo `layrz:"subform"`
	}

	form := &PersonForm{
		UserInfo: &UserInfo{
			FirstName: nil,
		},
	}

	errs := Validate(form)

	// Should have userInfo.firstName (both camelCased)
	if _, ok := errs["userInfo.firstName"]; !ok {
		t.Fatalf("expected userInfo.firstName key, got keys: %v", errs.Keys())
	}
}

// TestListIndexCamelCasing tests that numeric list indices are NOT camelCased.
func TestListIndexCamelCasing(t *testing.T) {
	type Item struct {
		ItemName *string `layrz:"char,required"`
	}

	type Form struct {
		ItemList []*Item `layrz:"subform_list"`
	}

	form := &Form{
		ItemList: []*Item{
			{ItemName: nil},
		},
	}

	errs := Validate(form)

	// Should have itemList.0.itemName (index "0" is NOT camelCased, but names are)
	if _, ok := errs["itemList.0.itemName"]; !ok {
		t.Fatalf("expected itemList.0.itemName key, got keys: %v", errs.Keys())
	}
}

// TestMultipleErrors per subform field tests that multiple validation errors are accumulated.
func TestMultipleErrorsPerField(t *testing.T) {
	type Address struct {
		StreetName *string `layrz:"char,required,min_length=5,max_length=20"`
	}

	type Form struct {
		Address *Address `layrz:"subform"`
	}

	form := &Form{
		Address: &Address{
			StreetName: Ptr("a"), // Too short
		},
	}

	errs := Validate(form)

	addrErrs, ok := errs["address.streetName"]
	if !ok {
		t.Fatal("missing address.streetName")
	}
	if len(addrErrs) != 1 || addrErrs[0].Code != testCodeMinLength {
		t.Errorf("expected 1 minLength error, got %v", addrErrs)
	}
}
