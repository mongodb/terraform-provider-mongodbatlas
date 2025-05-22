package schema

type CodeStatement struct {
	Code    string
	Imports []string
}

func GroupCodeStatements(stmts []CodeStatement, grouping func([]string) string) CodeStatement {
	listOfCode := []string{}
	imports := []string{}
	for i := range stmts {
		listOfCode = append(listOfCode, stmts[i].Code)
		imports = append(imports, stmts[i].Imports...)
	}
	resultCode := grouping(listOfCode)
	return CodeStatement{
		Code:    resultCode,
		Imports: imports,
	}
}
