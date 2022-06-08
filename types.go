package main

import "github.com/hasura/go-graphql-client"

type JobsPublic []struct {
	ID         graphql.ID
	Title      graphql.String
	Company    Company
	MinSalary  graphql.Int
	MaxSalary  graphql.Int
	Field      graphql.String
	Experience graphql.String
}

type Company struct {
	ID   graphql.ID
	Name graphql.String
}
