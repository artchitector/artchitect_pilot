package entropy

type Entropy struct {
}

/*
	Artchitect asks "select one element from set, i have total 100 elements.
	Entropy replies: "take element 31" (calculated with the lightnoise-entropy)
*/

func (e *Entropy) GetAnswer(totalElements uint) (uint, error) {
	return 0, nil
}
