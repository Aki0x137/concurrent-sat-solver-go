package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"maps"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Aki0x137/concurrent-sat-solver-go/set"
)

type Literal int
type Clause []Literal
type Formula []Clause

type Assignment map[Literal]bool

func dpll(formula Formula, assignment Assignment) (bool, Assignment) {
	if len(formula) == 0 {
		return false, assignment
	}

	if !checkClauseValidity(formula) {
		return false, assignment
	}

	if isSatisfied(formula, assignment) {
		return true, assignment
	}

	newFormula, assignment := unitPropagate(formula, assignment)

	newFormula, assignment = pureLiteralAssignment(newFormula, assignment)

	if isSatisfied(newFormula, assignment) {
		return true, assignment
	}

	if !checkClauseValidity(formula) {
		return false, assignment
	}

	selectedLiteral, err := selectLiteral(newFormula)
	if err != nil {
		return false, assignment
	}

	var simplifiedFormula Formula
	newAssignment := maps.Clone(assignment)
	newAssignment[selectedLiteral] = true
	for _, clause := range newFormula {
		if !slices.Contains(clause, selectedLiteral) {
			var updatedClause Clause
			if index := slices.Index(clause, -selectedLiteral); index >= 0 {
				updatedClause = slices.Delete(clause, index, index+1)
			} else {
				updatedClause = slices.Clone(clause)
			}

			simplifiedFormula = append(simplifiedFormula, updatedClause)
		}
	}

	result, finalAsassignment := dpll(simplifiedFormula, newAssignment)
	if result {
		return result, finalAsassignment
	}

	simplifiedFormula = make(Formula, 0)
	newAssignment = maps.Clone(assignment)
	newAssignment[selectedLiteral] = false
	for _, clause := range newFormula {
		if !slices.Contains(clause, -selectedLiteral) {
			var updatedClause Clause
			if index := slices.Index(clause, selectedLiteral); index >= 0 {
				updatedClause = slices.Delete(clause, index, index+1)
			} else {
				updatedClause = slices.Clone(clause)
			}

			simplifiedFormula = append(simplifiedFormula, updatedClause)
		}
	}
	return dpll(simplifiedFormula, newAssignment)
}

func checkClauseValidity(formula Formula) bool {
	for _, clause := range formula {
		if len(clause) == 0 {
			return false
		}
	}
	return true
}

func selectLiteral(formula Formula) (Literal, error) {
	for _, clause := range formula {
		for _, lit := range clause {
			return lit, nil
		}
	}
	return 0, errors.New("no literal found")
}

// isSatisified checks if all clauses are satisfied with the current assignment
func isSatisfied(formula Formula, assignment Assignment) bool {
	for _, clause := range formula {
		satisfied := false
		for _, literal := range clause {
			absVal := math.Abs(float64(literal))
			if val, ok := assignment[Literal(absVal)]; ok {
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
	for {
		var unitClauses []Clause
		for _, clause := range formula {
			if len(clause) == 1 {
				unitClauses = append(unitClauses, clause)
			}
		}

		if len(unitClauses) == 0 {
			break
		}

		for _, clause := range unitClauses {
			literal := clause[0]
			absVal := math.Abs(float64(literal))
			assignment[Literal(absVal)] = literal > 0
			for _, c := range updatedFormula {
				if !slices.Contains(c, literal) {
					updatedFormula = append(updatedFormula, c)
				}
				if index := slices.Index(c, -literal); index >= 0 {
					c = slices.Delete(c, index, index+1)
					updatedFormula = append(updatedFormula, c)
				}

			}
		}
	}

	return updatedFormula, assignment
}

/*
pureLiteralAssignment checks for pure literals and updates formula.

If a propositional variable occurs with only one polarity in the formula, it is called pure. A pure literal can always be assigned in a way that makes all clauses containing it true. Thus, when it is assigned in such a way, these clauses do not constrain the search anymore, and can be deleted.
*/
func pureLiteralAssignment(formula Formula, assignment Assignment) (Formula, Assignment) {
	updatedFormula := slices.Clone(formula)
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
		absVal := math.Abs(float64(literal))
		assignment[Literal(absVal)] = literal > 0

		for _, clause := range updatedFormula {
			if index := slices.Index(clause, literal); index >= 0 {
				updatedFormula = slices.Delete(updatedFormula, index, index+1)
			}
		}
	}

	return updatedFormula, assignment
}

func main() {
	file, err := os.Open("input.csv")

	if err != nil {
		log.Fatal("Error opening input file.\n Exiting...")
		return
	}
	defer file.Close()

	var formula Formula

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		literals := strings.Split(row, ",")
		var clause Clause

		for _, lit_str := range literals {
			literal, err := strconv.Atoi(lit_str)
			if err != nil {
				log.Fatal("Error while reading input.\nExiting...")
			}
			clause = append(clause, Literal(literal))
		}

		formula = append(formula, clause)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading input. \nExiting...")
	}

	var assignments Assignment
	sat, final_assignments := dpll(formula, assignments)

	if sat {
		fmt.Printf("The formula is satisfiable!\nAssignments %v", final_assignments)
	} else {
		fmt.Print("The formula can't be satisfied")
	}
}
