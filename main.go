package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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

	// Unit Propagation
	for _, clause := range formula {
		if len(clause) == 1 {
			literal := clause[0]
			value := (literal > 0)
			literal_abs := Literal(math.Abs(float64(literal))) // literal without negation, if any
			assignment[literal_abs] = value

			new_clauses := propagate(formula, literal)

			return dpll(new_clauses, assignment)
		}
	}

	return false, assignment
}

func checkClauseValidity(formula Formula) bool {
	for _, clause := range formula {
		if len(clause) == 0 {
			return false
		}
	}
	return true
}

func get_unassigned_literals(formula Formula, assignments Assignment) []Literal {

}

func propagate(formula Formula, literal Literal) Formula {

}

func select_literal(formula Formula) (error, Literal) {
	for _, clause := range formula {
		for _, lit := range clause {
			return nil, lit
		}
	}
	return errors.New("no literal found"), 0
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
