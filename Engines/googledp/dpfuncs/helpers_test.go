package dpfuncs

import (
	"testing"
)

// func isCorrectBudgetAndAllSubMeasHasBudget(query entities.Query, topLevelBudget entities.Budget, privacyNotation string) (bool, error) {
// Checks that either the top level budget is equal to all sub level budgets -> true
// Checks that either all sublevel budgets exist or no sub level budgets exist -> true
// If some sublevels budgets exist -> false
// If only toplevelbudget exist -> true
// If budget not according to privacynotation -> error

/*
	TESTS
	func isCorrectBudgetAndAllSubMeasHasBudget

	ID     			ExpectedRes		Tests
	TestHelper1		true, nil		only toplevel budget supplied
	TestHelper2		true, nil		all sublevel budgets supplied equal toplevel
	TestHelper3		false, nil		all sublevel budgets supplied not equal toplevel
	TestHelper4		false, nil		some sublevel budgets supplied
	TestHelper5		false, err		wrong budget schema according to PureDP
	TestHelper6     false, err		wrong budget schema according to ApproxDP


*/

func TestHelper1(t *testing.T) {

}

func TestHelper2(t *testing.T) {

}
func TestHelper3(t *testing.T) {

}
func TestHelper4(t *testing.T) {

}
func TestHelper5(t *testing.T) {

}
func TestHelper6(t *testing.T) {

}

// func correctBudget(budget entities.Budget, privacyNotation string) (float64, float64, error) {
// checks if PureDP that only epsilon is present in budget, returns epsilon, delta = 0, nil
// checks if ApproxDP that both epsilon and delta is present in budget, returns epsilon, delta, nil
// otherwise gives error

/*
	TESTS
	func correctBudget

	ID     			ExpectedRes		Tests
	TestHelper7		ep, 0, nil		PureDP ok
	TestHelper8		0, 0, err		PureDP not ok
	TestHelper9		ep, de, nil		ApproxDP ok
	TestHelper10	0, 0, err		ApproxDP not ok

*/

// func isValidShape(input [][]string) bool {
// Checks that a subquery is only in below formats
// Filter -> Bin -> measurement
// Filter -> measurment
// Bin -> measurment
// measurment
// otherwise returns false

/*
	TESTS
	func isValidShape


	ID     			ExpectedRes		Tests
	TestHelper11	true			measurment
	TestHelper12	true			filter -> measurement
	TestHelper13	false			measurement -> filter
	TestHelper14	true			filter -> bin -> measurement
	TestHelper15	true			bin -> measurement
	TestHelper16	false			measurement -> bin

*/

// func splitIntoSubQueries(query entities.Query) ([][]entities.QueryStep, error) {
// splits the querysteps into subqueries
// splits on measurments
// returns error if the shape is incorrect see isValidShape

/*
	TESTS
	func splitIntoSubQueries


	ID     			ExpectedRes		Tests
	TestHelper17	true			measurment | measurment
	TestHelper18	true			filter -> measurement | measurement
	TestHelper19	false			measurement -> filter | measurment
	TestHelper20	false			filter -> bin -> measurement | bin
	TestHelper21	false			bin -> measurement | bin
	TestHelper22	false			measurement -> bin | filter
	TestHelper23	true			filter -> bin -> measurement | filter -> measurement | measurement | bin -> measurement

*/

// func getIndexAndTypeFromSchema(columns []entities.Column, columnName string) (int, entities.ColType, error) {
// returns the index of the column, its columnType
// if column does not exist returns error

/*
	TESTS
	func getIndexAndTypeFromSchema


	ID     			ExpectedRes			Tests
	TestHelper24	index, type, nil	finds the column and index returned is corret
	TestHelper25	-1, nil, err		does not find the column

*/

// func checkIfNumber(name string) bool {
// checks if the provided string is name == Int or name == Double

/*
	TESTS
	func checkIfNumber


	ID     			ExpectedRes		Tests
	TestHelper26	true			name = Int, true
	TestHelper27	true			name = Double, true
	TestHelper28	false			name = xyz, false
	TestHelper29	false			name = Text, false
*/

// func getNoiseGenerator(name string) noise.Noise
// returns the googledp noise generator
// laplace generator if name == "Laplace"
// gaussian generator if name == "Gaussian"
