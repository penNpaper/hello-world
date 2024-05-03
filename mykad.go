package mykad

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

const minAge = 12
const maxAge = 112 // Based off guinness world record.

// NewMyKAD returns a new malaysian identity from a provided NRIC number. The NRIC number can be in formatted or
// unformatted (no dash) format.
func NewMyKAD(nric string) (*MyKAD, error) {
	s := decodeNRIC(nric)
	if len(s) != 4 {
		return nil, errors.New("invalid mykad number")
	}

	// Parse date of birth.
	dob, err := time.Parse("060102", s[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %v", err)
	}

	// Parse place of birth.
	pob, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing location: %v", err)
	}

	var c CitizenType
	var pb PlaceOfBirth
	pobc := placeOfBirthCode(pob)
	if !pobc.valid() {
		return nil, errors.New("invalid place of birth")
	}

	re := regions[pobc]
	if pobc.isMalaysianCode() {
		c = CitizenTypeMalaysian
		pb = PlaceOfBirth{
			Country:  "Malaysia",
			Province: re,
		}
	} else if pobc.isForeignerCode() {
		c = CitizenTypeForeigner
		pb = PlaceOfBirth{
			Country:  re,
			Province: "",
		}
	}

	// Parse gender.
	gender := GenderMale
	g, err := strconv.Atoi(s[3])
	if err != nil {
		return nil, fmt.Errorf("error parsing gender: %v", err)
	}

	if g%2 == 0 {
		gender = GenderFemale
	}

	return &MyKAD{
		NRIC:         nric,
		DateOfBirth:  dob,
		PlaceOfBirth: pb,
		CitizenType:  c,
		Gender:       gender,
	}, nil
}

// decodeNRIC is a utility function that returns a split NRIC string.
func decodeNRIC(nric string) []string {
	// Try without dashes.
	if len(nric) == 12 {
		return []string{nric[0:6], nric[6:8], nric[8:11], nric[11:]}
	}

	// Try with dashes.
	r := regexp.MustCompile(`^(\d{6})-?(\d{2})-?(\d{3})(\d{1})$`)
	return r.FindStringSubmatch(nric)[1:]
}

// Generate will return a new random MyKAD number.
func Generate() string {
	// Generate a random date for the year component.
	n := time.Now()
	rand.Seed(n.UnixNano())

	e := n.AddDate(-minAge, 0, 0).Unix()
	s := n.AddDate(-maxAge, 0, 0).Unix()

	sec := rand.Int63n(e-s) + s
	ds := time.Unix(sec, 0).Format("060102")

	// Generate a random place.
	var p placeOfBirthCode
	for !p.valid() {
		p = placeOfBirthCode(rand.Intn(99))
	}

	// Generate a special number.
	sn := rand.Intn(9999)

	return fmt.Sprintf("%v-%02d-%04d", ds, p, sn)
}

// Validate will validate a NRIC number.
func Validate(nric string) error {
	_, err := NewMyKAD(nric)
	return err
}
