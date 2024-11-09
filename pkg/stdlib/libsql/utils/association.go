package utils

type QueryAssociation struct {
	TableName       interface{}
	FilterCondition string
	FilterArgs      []interface{}
}

type QueryAssociationBuilder struct {
	queryAssociation QueryAssociation
}

func NewQueryAssociationBuilder() *QueryAssociationBuilder {
	return &QueryAssociationBuilder{
		queryAssociation: QueryAssociation{},
	}
}

func (b *QueryAssociationBuilder) FilterCondition(filterCondition string, filterArgs ...interface{}) *QueryAssociationBuilder {
	b.queryAssociation.FilterCondition = filterCondition
	b.queryAssociation.FilterArgs = filterArgs

	return b
}

func (b *QueryAssociationBuilder) Build(tableName interface{}) QueryAssociation {
	b.queryAssociation.TableName = tableName
	return b.queryAssociation
}

func (b *QueryAssociationBuilder) BuildList(tableNames []interface{}) []QueryAssociation {
	var result []QueryAssociation
	for _, tableName := range tableNames {
		result = append(result, NewQueryAssociationBuilder().
			FilterCondition(b.queryAssociation.FilterCondition, b.queryAssociation.FilterArgs...).
			Build(tableName))
	}

	return result
}
