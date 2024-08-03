package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Aki0x137/concurrent-sat-solver-go/dpll"
)

// type Literal int
// type Clause []Literal
// type Formula []Clause

// type Assignment map[Literal]bool

// func checkClauseValidity(formula Formula) bool {
// 	for _, clause := range formula {
// 		if len(clause) == 0 {
// 			return false
// 		}
// 	}
// 	return true
// }

// func selectLiteral(formula Formula, assignment Assignment) (Literal, error) {
// 	for _, clause := range formula {
// 		for _, literal := range clause {
// 			absLiteral := math.Abs(float64(literal))
// 			if _, ok := assignment[Literal(absLiteral)]; !ok {
// 				return literal, nil
// 			}
// 		}
// 	}
// 	return 0, errors.New("no literal found")
// }

// // isSatisified checks if all clauses are satisfied with the current assignment
// func isSatisfied(formula Formula, assignment Assignment) bool {
// 	for _, clause := range formula {
// 		satisfied := false
// 		for _, literal := range clause {
// 			absLiteral := math.Abs(float64(literal))
// 			if val, ok := assignment[Literal(absLiteral)]; ok {
// 				if (literal > 0 && val) || (literal < 0 && !val) {
// 					satisfied = true
// 					break
// 				}
// 			}
// 		}
// 		if !satisfied {
// 			return false
// 		}
// 	}

// 	return true
// }

// // unitPropagate performs unit propagation on formula, based on curent assignments
// func unitPropagate(formula Formula, assignment Assignment) (Formula, Assignment) {
// 	updatedFormula := slices.Clone(formula)
// 	updatedAssignment := maps.Clone(assignment)
// 	for {
// 		var unitClauses []Clause
// 		for _, clause := range updatedFormula {
// 			if len(clause) == 1 {
// 				unitClauses = append(unitClauses, clause)
// 			}
// 		}

// 		if len(unitClauses) == 0 {
// 			break
// 		}

// 		for _, clause := range unitClauses {
// 			literal := clause[0]
// 			absLiteral := math.Abs(float64(literal))
// 			updatedAssignment[Literal(absLiteral)] = literal > 0

// 			var filteredFormula Formula
// 			for _, c := range updatedFormula {
// 				if !slices.Contains(c, literal) {
// 					filteredFormula = append(filteredFormula, c)
// 				}
// 			}

// 			var simplifiedFormula Formula
// 			for _, c := range filteredFormula {
// 				updatedClause := slices.Clone(c)
// 				if index := slices.Index(updatedClause, -literal); index >= 0 {
// 					updatedClause = slices.Delete(updatedClause, index, index+1)
// 				}
// 				simplifiedFormula = append(simplifiedFormula, updatedClause)
// 			}
// 			updatedFormula = simplifiedFormula
// 		}
// 	}

// 	return updatedFormula, updatedAssignment
// }

// /*
// pureLiteralAssignment checks for pure literals and updates formula.

// If a propositional variable occurs with only one polarity in the formula, it is called pure. A pure literal can always be assigned in a way that makes all clauses containing it true. Thus, when it is assigned in such a way, these clauses do not constrain the search anymore, and can be deleted.
// */
// func pureLiteralAssignment(formula Formula, assignment Assignment) (Formula, Assignment) {
// 	updatedFormula := slices.Clone(formula)
// 	updatedAssignment := maps.Clone(assignment)

// 	allLiteralsSet := set.NewSet[Literal]()
// 	for _, clauses := range formula {
// 		for _, literal := range clauses {
// 			allLiteralsSet.Add(literal)
// 		}
// 	}

// 	allLiterals := allLiteralsSet.Values()
// 	pureLiterals := set.NewSet[Literal]()
// 	for _, literal := range allLiterals {
// 		if !slices.Contains(allLiterals, -literal) {
// 			pureLiterals.Add(literal)
// 		}
// 	}

// 	for _, literal := range pureLiterals.Values() {
// 		absLiteral := math.Abs(float64(literal))
// 		updatedAssignment[Literal(absLiteral)] = literal > 0

// 		var filteredFormula Formula
// 		for _, clause := range updatedFormula {
// 			if index := slices.Index(clause, literal); index == -1 {
// 				filteredFormula = append(filteredFormula, clause)
// 			}
// 		}
// 		updatedFormula = filteredFormula
// 	}

// 	return updatedFormula, updatedAssignment
// }

// func simplifyFormula(formula Formula, literal Literal) Formula {
// 	var simplifiedFormula Formula
// 	for _, clause := range formula {
// 		if !slices.Contains(clause, literal) {
// 			updatedClause := slices.Clone(clause)
// 			if index := slices.Index(updatedClause, -literal); index >= 0 {
// 				updatedClause = slices.Delete(updatedClause, index, index+1)
// 			}
// 			simplifiedFormula = append(simplifiedFormula, updatedClause)
// 		}
// 	}
// 	return simplifiedFormula
// }

// func dpllSolve(formula Formula, assignment Assignment) (bool, Assignment) {
// 	if len(formula) == 0 {
// 		return false, assignment
// 	}

// 	if !checkClauseValidity(formula) {
// 		return false, assignment
// 	}

// 	if isSatisfied(formula, assignment) {
// 		return true, assignment
// 	}

// 	newFormula, newAssignment := unitPropagate(formula, assignment)

// 	newFormula, newAssignment = pureLiteralAssignment(newFormula, newAssignment)

// 	if isSatisfied(newFormula, newAssignment) {
// 		return true, newAssignment
// 	}

// 	if !checkClauseValidity(formula) {
// 		return false, assignment
// 	}

// 	selectedLiteral, err := selectLiteral(newFormula, newAssignment)
// 	if err != nil {
// 		return false, assignment
// 	}

// 	assignment1 := maps.Clone(newAssignment)
// 	assignment1[selectedLiteral] = true

// 	simplifiedFormula := simplifyFormula(newFormula, selectedLiteral)

// 	result, finalAssignment := dpllSolve(simplifiedFormula, assignment1)
// 	if result {
// 		return result, finalAssignment
// 	}

// 	assignment2 := maps.Clone(newAssignment)
// 	assignment2[selectedLiteral] = false

// 	simplifiedFormula = simplifyFormula(newFormula, -selectedLiteral)

// 	result, finalAssignment = dpllSolve(simplifiedFormula, assignment2)
// 	return result, finalAssignment
// }

// func recursiveSolution(formula Formula, assignment Assignment, ch chan Assignment) {
// 	if len(formula) == 0 {
// 		return
// 	}

// 	if !checkClauseValidity(formula) {
// 		return
// 	}

// 	if isSatisfied(formula, assignment) {
// 		ch <- assignment
// 		close(ch)
// 		return
// 	}

// 	newFormula, newAssignment := unitPropagate(formula, assignment)

// 	newFormula, newAssignment = pureLiteralAssignment(newFormula, newAssignment)

// 	if isSatisfied(newFormula, newAssignment) {
// 		ch <- newAssignment
// 		close(ch)
// 		return
// 	}

// 	if !checkClauseValidity(formula) {
// 		return
// 	}

// 	selectedLiteral, err := selectLiteral(newFormula, newAssignment)
// 	if err != nil {
// 		return
// 	}

// 	assignment1 := maps.Clone(newAssignment)
// 	assignment1[selectedLiteral] = true

// 	simplifiedFormula := simplifyFormula(newFormula, selectedLiteral)

// 	result, finalAssignment := dpllSolve(simplifiedFormula, assignment1)
// 	if result {
// 		ch <- finalAssignment
// 		close(ch)
// 		return
// 	}

// 	assignment2 := maps.Clone(newAssignment)
// 	assignment2[selectedLiteral] = false

// 	simplifiedFormula = simplifyFormula(newFormula, -selectedLiteral)

// 	result, finalAssignment = dpllSolve(simplifiedFormula, assignment2)
// 	if result {
// 		ch <- finalAssignment
// 		close(ch)
// 		return
// 	}
// }

func main() {
	file, err := os.Open("input.csv")

	if err != nil {
		log.Fatal("Error opening input file.\n Exiting...")
		return
	}
	defer file.Close()

	var formula dpll.Formula

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		literals := strings.Split(row, ",")
		var clause dpll.Clause

		for _, lit_str := range literals {
			literal, err := strconv.Atoi(lit_str)
			if err != nil {
				log.Fatal("Error while reading input.\nExiting...")
			}
			clause = append(clause, dpll.Literal(literal))
		}

		formula = append(formula, clause)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading input. \nExiting...")
	}

	sat, final_assignments := dpll.Solve(formula, make(dpll.Assignment))

	if sat {
		fmt.Printf("The formula is satisfiable!\nAssignments %v\n", final_assignments)
	} else {
		fmt.Println("The formula can't be satisfied")
	}
}
