package dpll

import (
	"errors"
	"maps"
	"math"
	"slices"
	"sync"

	"github.com/Aki0x137/concurrent-sat-solver-go/set"
)

type Literal int
type Clause []Literal
type Formula []Clause

type Assignment map[Literal]bool

// checkClauseValidity validates the formula
func checkClauseValidity(formula Formula) bool {
	for _, clause := range formula {
		if len(clause) == 0 {
			return false
		}
	}
	return true
}

// selectLiteral chooses an unassigned literal in the formula. If no literal is found, it return an error
func selectLiteral(formula Formula, assignment Assignment) (Literal, error) {
	for _, clause := range formula {
		for _, literal := range clause {
			absLiteral := math.Abs(float64(literal))
			if _, ok := assignment[Literal(absLiteral)]; !ok {
				return literal, nil
			}
		}
	}
	return 0, errors.New("no literal found")
}

// isSatisfied checks if all clauses are satisfied with the current assignment
func isSatisfied(formula Formula, assignment Assignment) bool {
	for _, clause := range formula {
		satisfied := false
		for _, literal := range clause {
			absLiteral := math.Abs(float64(literal))
			if val, ok := assignment[Literal(absLiteral)]; ok {
				if (literal > 0 && val) || (literal < 0 && !val) {
					satisfied = true
					break
				}
			}
		}
		if !satisfied {
			return false
		}
	}

	return true
}

// unitPropagate performs unit propagation on formula, based on curent assignments
func unitPropagate(formula Formula, assignment Assignment) (Formula, Assignment) {
	updatedFormula := slices.Clone(formula)
	updatedAssignment := maps.Clone(assignment)
	for {
		var unitClauses []Clause
		for _, clause := range updatedFormula {
			if len(clause) == 1 {
				unitClauses = append(unitClauses, clause)
			}
		}

		if len(unitClauses) == 0 {
			break
		}

		for _, clause := range unitClauses {
			literal := clause[0]
			absLiteral := math.Abs(float64(literal))
			updatedAssignment[Literal(absLiteral)] = literal > 0

			var filteredFormula Formula
			for _, c := range updatedFormula {
				if !slices.Contains(c, literal) {
					filteredFormula = append(filteredFormula, c)
				}
			}

			var simplifiedFormula Formula
			for _, c := range filteredFormula {
				updatedClause := slices.Clone(c)
				if index := slices.Index(updatedClause, -literal); index >= 0 {
					updatedClause = slices.Delete(updatedClause, index, index+1)
				}
				simplifiedFormula = append(simplifiedFormula, updatedClause)
			}
			updatedFormula = simplifiedFormula
		}
	}

	return updatedFormula, updatedAssignment
}

/*
pureLiteralAssignment checks for pure literals and updates formula.

If a propositional variable occurs with only one polarity in the formula, it is called pure. A pure literal can always be assigned in a way that makes all clauses containing it true. Thus, when it is assigned in such a way, these clauses do not constrain the search anymore, and can be deleted.
*/
func pureLiteralAssignment(formula Formula, assignment Assignment) (Formula, Assignment) {
	updatedFormula := slices.Clone(formula)
	updatedAssignment := maps.Clone(assignment)

	allLiteralsSet := set.NewSet[Literal]()
	for _, clauses := range formula {
		for _, literal := range clauses {
			allLiteralsSet.Add(literal)
		}
	}

	allLiterals := allLiteralsSet.Values()
	pureLiterals := set.NewSet[Literal]()
	for _, literal := range allLiterals {
		if !slices.Contains(allLiterals, -literal) {
			pureLiterals.Add(literal)
		}
	}

	for _, literal := range pureLiterals.Values() {
		absLiteral := math.Abs(float64(literal))
		updatedAssignment[Literal(absLiteral)] = literal > 0

		var filteredFormula Formula
		for _, clause := range updatedFormula {
			if index := slices.Index(clause, literal); index == -1 {
				filteredFormula = append(filteredFormula, clause)
			}
		}
		updatedFormula = filteredFormula
	}

	return updatedFormula, updatedAssignment
}

// simplifyFormula removes truthy clauses and removes redundant literals after an assignment
func simplifyFormula(formula Formula, literal Literal) Formula {
	var simplifiedFormula Formula
	for _, clause := range formula {
		if !slices.Contains(clause, literal) {
			updatedClause := slices.Clone(clause)
			if index := slices.Index(updatedClause, -literal); index >= 0 {
				updatedClause = slices.Delete(updatedClause, index, index+1)
			}
			simplifiedFormula = append(simplifiedFormula, updatedClause)
		}
	}
	return simplifiedFormula
}

func Solve(formula Formula, assignment Assignment) (bool, Assignment) {
	if len(formula) == 0 {
		return false, assignment
	}

	if !checkClauseValidity(formula) {
		return false, assignment
	}

	if isSatisfied(formula, assignment) {
		return true, assignment
	}

	newFormula, newAssignment := unitPropagate(formula, assignment)

	newFormula, newAssignment = pureLiteralAssignment(newFormula, newAssignment)

	if isSatisfied(newFormula, newAssignment) {
		return true, newAssignment
	}

	if !checkClauseValidity(formula) {
		return false, assignment
	}

	selectedLiteral, err := selectLiteral(newFormula, newAssignment)
	if err != nil {
		return false, assignment
	}

	var wg sync.WaitGroup
	ch := make(chan Assignment, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		assignment1 := maps.Clone(newAssignment)
		assignment1[selectedLiteral] = true
		simplifiedFormula := simplifyFormula(newFormula, selectedLiteral)
		recursiveSolution(simplifiedFormula, assignment1, ch)
	}()

	go func() {
		defer wg.Done()
		assignment2 := maps.Clone(newAssignment)
		assignment2[selectedLiteral] = false
		simplifiedFormula := simplifyFormula(newFormula, -selectedLiteral)
		recursiveSolution(simplifiedFormula, assignment2, ch)
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for finalAssignment := range ch {
		if isSatisfied(formula, finalAssignment) {
			return true, finalAssignment
		}
	}

	return false, assignment
}

func recursiveSolution(formula Formula, assignment Assignment, ch chan Assignment) {
	if len(formula) == 0 {
		return
	}

	if !checkClauseValidity(formula) {
		return
	}

	if isSatisfied(formula, assignment) {
		ch <- assignment
		return
	}

	newFormula, newAssignment := unitPropagate(formula, assignment)

	newFormula, newAssignment = pureLiteralAssignment(newFormula, newAssignment)

	if isSatisfied(newFormula, newAssignment) {
		ch <- newAssignment
		return
	}

	if !checkClauseValidity(formula) {
		return
	}

	selectedLiteral, err := selectLiteral(newFormula, newAssignment)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	subCh := make(chan Assignment, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		assignment1 := maps.Clone(newAssignment)
		assignment1[selectedLiteral] = true
		simplifiedFormula := simplifyFormula(newFormula, selectedLiteral)
		recursiveSolution(simplifiedFormula, assignment1, subCh)
	}()

	go func() {
		defer wg.Done()
		assignment2 := maps.Clone(newAssignment)
		assignment2[selectedLiteral] = false
		simplifiedFormula := simplifyFormula(newFormula, -selectedLiteral)
		recursiveSolution(simplifiedFormula, assignment2, subCh)
	}()

	go func() {
		wg.Wait()
		close(subCh)
	}()

	for finalAssignment := range subCh {
		if isSatisfied(formula, finalAssignment) {
			ch <- finalAssignment
			return
		}
	}
}
