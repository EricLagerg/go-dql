package godql

/*
	Package godql provides an easy-to-use interface for generating DQL
	statements.

	godql IS NOT AN ORM.

	For example,
		"SELECT DISTINCT f1 WHERE stateabbr='XX' AND distr=12 LIMIT 5"
	would be written as:
		query := new(godql.Query)
		query.SelectDistinct("fi").
		Where("stateabbr", Equals, "XX").
		Where("distr", Equals, 12).
		Limit(5)

	Brief notes:
		SelectDistinct() has a higher priority than Select(), meaning
		it'll override a Select() statement.

		And() and Or() should only be used to group two or more of
		their respective statements.

*/
