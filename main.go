package main

import (
	. "./godql"
	"fmt"
)

const test = "SELECT RNCCalcParty,count(RNCCalcParty) WHERE stateabbreviation='IL' AND CongressionalDistrict=12 GROUP BY RNCCalcParty LIMIT 200"

func main() {
	query := new(Query)
	query.Select([]string{"RNCCalcParty"}).
		Count("RNCCalcParty").
		Where("stateabbreviation", Equals, "IL").
		Where("CongressionalDistrict", Equals, 12).
		GroupBy("RNCCalcParty").
		Limit(200)

	foo := query.String()
	if foo != test {
		fmt.Println(foo)
		fmt.Println(test)
	}
}
