package dpll

import (
	"errors"
	"fmt"
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

type Result struct {
	mu         *sync.RWMutex
	satisfied  bool
	assignment Assignment
}

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

// isSatisified checks if all clauses are satisfied with the current assignment
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

	ans := Result{&sync.RWMutex{}, false, Assignment{}}
	var wg sync.WaitGroup

	assignment1 := maps.Clone(newAssignment)
	assignment1[selectedLiteral] = true

	simplifiedFormula := simplifyFormula(newFormula, selectedLiteral)

	wg.Add(1)
	go recursiveSolution(simplifiedFormula, assignment1, &wg, &ans)

	assignment2 := maps.Clone(newAssignment)
	assignment2[selectedLiteral] = false

	simplifiedFormula = simplifyFormula(newFormula, -selectedLiteral)

	wg.Add(1)
	go recursiveSolution(simplifiedFormula, assignment2, &wg, &ans)

	wg.Wait()

	fmt.Println("After wait")

	if ans.satisfied {
		return true, ans.assignment
	} else {
		return false, assignment
	}
}

func recursiveSolution(formula Formula, assignment Assignment, wg *sync.WaitGroup, ans *Result) {
	defer wg.Done()

	fmt.Println("INside recursive", assignment)

	ans.mu.RLock()
	if ans.satisfied {
		return
	}
	ans.mu.RUnlock()

	if len(formula) == 0 {
		return
	}

	if !checkClauseValidity(formula) {
		return
	}

	if isSatisfied(formula, assignment) {
		ans.mu.Lock()
		ans.assignment = assignment
		ans.satisfied = true
		ans.mu.Unlock()
		return
	}

	newFormula, newAssignment := unitPropagate(formula, assignment)

	newFormula, newAssignment = pureLiteralAssignment(newFormula, newAssignment)

	if isSatisfied(newFormula, newAssignment) {
		ans.mu.Lock()
		ans.assignment = newAssignment
		ans.satisfied = true
		ans.mu.Unlock()
		return
	}

	if !checkClauseValidity(formula) {
		return
	}

	selectedLiteral, err := selectLiteral(newFormula, newAssignment)
	if err != nil {
		return
	}

	assignment1 := maps.Clone(newAssignment)
	assignment1[selectedLiteral] = true

	simplifiedFormula := simplifyFormula(newFormula, selectedLiteral)

	wg.Add(1)
	go recursiveSolution(simplifiedFormula, assignment1, wg, ans)

	assignment2 := maps.Clone(newAssignment)
	assignment2[selectedLiteral] = false

	simplifiedFormula = simplifyFormula(newFormula, -selectedLiteral)

	wg.Add(1)
	go recursiveSolution(simplifiedFormula, assignment2, wg, ans)
}
