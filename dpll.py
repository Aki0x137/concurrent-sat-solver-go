def dpll(clauses, assignment):
    """
    DPLL algorithm to determine satisfiability of a set of clauses.

    Parameters:
    - clauses: List of lists representing the CNF clauses.
    - assignment: Dictionary representing the current variable assignments.

    Returns:
    - (bool, dict): A tuple where the first element is True if the clauses are satisfiable, and
                    the second element is the assignment that satisfies the clauses.
    """

    # If all clauses are satisfied
    if is_satisfied(clauses, assignment):
        return True, assignment

    # If there is an empty clause
    if any(len(clause) == 0 for clause in clauses):
        return False, {}

    # Perform unit propagation
    clauses, assignment = unit_propagate(clauses, assignment)

    # Perform pure literal assignment
    clauses, assignment = pure_literal_assign(clauses, assignment)

    # If all clauses are satisfied after unit propagation and pure literal assignment
    if is_satisfied(clauses, assignment):
        return True, assignment

    # If there is an empty clause after unit propagation and pure literal assignment
    if any(len(clause) == 0 for clause in clauses):
        return False, {}

    # Choose an unassigned variable
    for clause in clauses:
        for literal in clause:
            var = abs(literal)
            if var not in assignment:
                break
        else:
            continue
        break

    # Recursively check with the variable set to True
    new_assignment = assignment.copy()
    new_assignment[var] = True
    simplified_clauses = [[l for l in c if l != -var] for c in clauses if var not in c]
    result, final_assignment = dpll(simplified_clauses, new_assignment)
    if result:
        return True, final_assignment

    # Recursively check with the variable set to False
    new_assignment = assignment.copy()
    new_assignment[var] = False
    simplified_clauses = [[l for l in c if l != var] for c in clauses if -var not in c]
    return dpll(simplified_clauses, new_assignment)

def is_satisfied(clauses, assignment):
    """Check if all clauses are satisfied with the current assignment."""
    for clause in clauses:
        satisfied = False
        for literal in clause:
            var = abs(literal)
            if var in assignment:
                val = assignment[var]
                if (literal > 0 and val) or (literal < 0 and not val):
                    satisfied = True
                    break
        if not satisfied:
            return False
    return True

def unit_propagate(clauses, assignment):
    """Perform unit propagation."""
    while True:
        unit_clauses = [clause for clause in clauses if len(clause) == 1]
        if not unit_clauses:
            break
        for clause in unit_clauses:
            literal = clause[0]
            var = abs(literal)
            val = literal > 0
            assignment[var] = val
            clauses = [c for c in clauses if literal not in c]
            clauses = [[l for l in c if l != -literal] for c in clauses]
    return clauses, assignment

def pure_literal_assign(clauses, assignment):
    """Assign pure literals."""
    all_literals = [literal for clause in clauses for literal in clause]
    pure_literals = set(literal for literal in all_literals if -literal not in all_literals)
    for literal in pure_literals:
        var = abs(literal)
        val = literal > 0
        assignment[var] = val
        clauses = [c for c in clauses if literal not in c]
        return clauses, assignment

# Test the DPLL algorithm with the given inputs
clauses1 = [[1, 2, 3], [-1, 2, -3], [1, -2, 3], [-1, -2, -3]]
clauses2 = [[1,2,3],[1,2,-3],[1,-2,3],[1,-2,-3],[-1,2,3],[-1,2,-3],[-1,-2,3],[-1,-2,-3]] # [[1, 2, -3], [-1, -2, 3], [-1, -2, -3], [1, 2, 3]]

print("Testing clauses1: ", clauses1)
result1, assignment1 = dpll(clauses1, {})
print("Satisfiable: ", result1)
print("Assignment: ", assignment1)

print("\nTesting clauses2: ", clauses2)
result2, assignment2 = dpll(clauses2, {})
print("Satisfiable: ", result2)
print("Assignment: ", assignment2)
