package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"
	"errors"
    "fmt"
	"github.com/vibovenkat123/pennyearned-server/graph/model"
	dbHelpers "github.com/vibovenkat123/pennyearned-server/internal/db/helpers"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user, err := dbHelpers.SignUp(input.Name, input.Username, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	data := &model.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
		Username:    user.Username,
		Password:    input.Password,
	}
	return data, nil
}

// CreateExpense is the resolver for the createExpense field.
func (r *mutationResolver) CreateExpense(ctx context.Context, input model.NewExpense) (*model.Expense, error) {
	expense, err := dbHelpers.NewExpense(input.OwnerID, input.Name, input.Spent)
	if err != nil {
		return nil, err
	}
	data := &model.Expense{
		ID:      expense.ID,
		OwnerID: expense.OwnerID,
		Name:    expense.Name,
		Spent:   expense.Spent,
	}
	return data, nil
}

// UpdateExpense is the resolver for the updateExpense field.
func (r *mutationResolver) UpdateExpense(ctx context.Context, input model.UpdateExpenseInput) (*model.Expense, error) {
	original, err := dbHelpers.GetExpenseById(input.ID)
	if err != nil {
		return nil, err
	}
	id := input.ID
	ownerId := input.OwnerID
	name := input.Name
	spent := input.Spent
	if ownerId == nil {
		ownerId = &original.OwnerID
	}
	if name == nil {
		name = &original.Name
	}
	if spent == nil {
		spent = &original.Spent
	}
	expense := dbHelpers.Expense{
		ID:      id,
		Name:    *name,
		Spent:   *spent,
		OwnerID: *ownerId,
	}
	expense, err = dbHelpers.UpdateExpense(expense)
	if err != nil {
		return nil, err
	}
	data := &model.Expense{
		ID:          expense.ID,
		OwnerID:     expense.OwnerID,
		Name:        expense.Name,
		Spent:       expense.Spent,
		DateCreated: expense.DateCreated,
		DateUpdated: expense.DateUpdated,
	}
	if data.ID != input.ID {
		return nil, errors.New("expense not found")
	}
	return data, nil
}

// DeleteExpense is the resolver for the deleteExpense field.
func (r *mutationResolver) DeleteExpense(ctx context.Context, input model.DeleteExpense) (*model.Expense, error) {
	expense, err := dbHelpers.DeleteExpense(input.ID)
	if err != nil {
		return nil, err
	}
	data := &model.Expense{
		ID:      expense.ID,
		OwnerID: expense.OwnerID,
		Name:    expense.Name,
		Spent:   expense.Spent,
	}
	if data.ID != input.ID {
		return nil, errors.New("expense not found")
	}
	return data, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	original, err := dbHelpers.SignIn(input.Email, input.Password)
	if err != nil {
        fmt.Println(original)
		return nil, err
	}
    var pass string
	name := input.Name
    email := input.NewEmail
    username := input.Username
    inputPass := input.NewPass
	if name == nil {
		name = &original.Name
	}
	if email == nil {
		email = &original.Email
	}
	if inputPass == nil {
        pass = original.Password
	} else {
        pass, err = dbHelpers.GenerateFromPassword(*inputPass, dbHelpers.P)
        if err != nil {
            return nil, err
        }
    }
    if username == nil {
        username = &original.Username
    }
	userInput := dbHelpers.User{
        Name: *name,
        Password: pass,
        Email: *email,
        Username: *username,
        ID: original.ID,
	}
	userInput, err = dbHelpers.UpdateUser(userInput)
	if err != nil {
		return nil, err
	}
	data := &model.User{
		ID:          userInput.ID,
		Name:        userInput.Name,
		DateCreated: userInput.DateCreated,
		DateUpdated: userInput.DateUpdated,
        Username: userInput.Username,
        Password: userInput.Password,
        Email: userInput.Email,
	}
	return data, nil
}

// Expenses is the resolver for the expenses field.
func (r *queryResolver) Expenses(ctx context.Context) ([]*model.Expense, error) {
	var expenses []*model.Expense
	data, err := dbHelpers.GetAllExpenses()
	if err != nil {
		return nil, err
	}
	for _, r := range data {
		expenses = append(expenses, &model.Expense{
			ID:          r.ID,
			OwnerID:     r.OwnerID,
			Name:        r.Name,
			Spent:       r.Spent,
			DateCreated: r.DateCreated,
			DateUpdated: r.DateUpdated,
		})
	}
	return expenses, nil
}

// Expense is the resolver for the expense field.
func (r *queryResolver) Expense(ctx context.Context, id string) (*model.Expense, error) {
	data, err := dbHelpers.GetExpenseById(id)
	if err != nil {
		return nil, err
	}
	expense := &model.Expense{
		ID:          data.ID,
		OwnerID:     data.OwnerID,
		Name:        data.Name,
		Spent:       data.Spent,
		DateCreated: data.DateCreated,
		DateUpdated: data.DateUpdated,
	}
	if expense.ID != id {
		return nil, errors.New("expense not found")
	}
	return expense, nil
}

// ExpensesHigherThan is the resolver for the expensesHigherThan field.
func (r *queryResolver) ExpensesHigherThan(ctx context.Context, spent int, ownerID string) ([]*model.Expense, error) {
	var expenses []*model.Expense
	data, err := dbHelpers.ExpensesHigherThan(spent, ownerID)
	if err != nil {
		return nil, err
	}
	for _, r := range data {
		expenses = append(expenses, &model.Expense{
			ID:          r.ID,
			OwnerID:     r.OwnerID,
			Name:        r.Name,
			Spent:       r.Spent,
			DateCreated: r.DateCreated,
			DateUpdated: r.DateUpdated,
		})
	}
	if len(expenses) <= 0 {
		return nil, errors.New("expenses not found")
	}
	return expenses, nil
}

// ExpensesLowerThan is the resolver for the expensesLowerThan field.
func (r *queryResolver) ExpensesLowerThan(ctx context.Context, spent int, ownerID string) ([]*model.Expense, error) {
	var expenses []*model.Expense
	data, err := dbHelpers.ExpensesLowerThan(spent, ownerID)
	if err != nil {
		return nil, err
	}
	for _, r := range data {
		expenses = append(expenses, &model.Expense{
			ID:          r.ID,
			OwnerID:     r.OwnerID,
			Name:        r.Name,
			Spent:       r.Spent,
			DateCreated: r.DateCreated,
			DateUpdated: r.DateUpdated,
		})
	}
	if len(expenses) <= 0 {
		return nil, errors.New("expenses not found")
	}
	return expenses, nil
}

// MostExpensiveExpense is the resolver for the mostExpensiveExpense field.
func (r *queryResolver) MostExpensiveExpense(ctx context.Context, ownerID string) (*model.Expense, error) {
	data, err := dbHelpers.MostExpensiveExpense(ownerID)
	if err != nil {
		return nil, err
	}
	expense := &model.Expense{
		ID:          data.ID,
		OwnerID:     data.OwnerID,
		Name:        data.Name,
		Spent:       data.Spent,
		DateCreated: data.DateCreated,
		DateUpdated: data.DateUpdated,
	}
	if expense.ID == "" {
		return nil, errors.New("There are no expenses")
	}
	return expense, nil
}

// LeastExpensiveExpense is the resolver for the leastExpensiveExpense field.
func (r *queryResolver) LeastExpensiveExpense(ctx context.Context, ownerID string) (*model.Expense, error) {
	data, err := dbHelpers.LeastExpensiveExpense(ownerID)
	if err != nil {
		return nil, err
	}
	expense := &model.Expense{
		ID:          data.ID,
		OwnerID:     data.OwnerID,
		Name:        data.Name,
		Spent:       data.Spent,
		DateCreated: data.DateCreated,
		DateUpdated: data.DateUpdated,
	}
	if expense.ID == "" {
		return nil, errors.New("There are no expenses")
	}
	return expense, nil
}

// ExpensesFromOwner is the resolver for the expensesFromOwner field.
func (r *queryResolver) ExpensesFromOwner(ctx context.Context, id string) ([]*model.Expense, error) {
	data, err := dbHelpers.GetExpensesByOwnerId(id)
	if err != nil {
		return nil, err
	}
	var expenses []*model.Expense
	for _, r := range data {
		if r.ID == "" {
			return nil, errors.New("There are no expenses")
		}
		expenses = append(expenses, &model.Expense{
			ID:          r.ID,
			OwnerID:     r.OwnerID,
			Name:        r.Name,
			Spent:       r.Spent,
			DateCreated: r.DateCreated,
			DateUpdated: r.DateUpdated,
		})
	}
	return expenses, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := dbHelpers.SignIn(email, password)
	if err != nil {
		return nil, err
	}
	data := &model.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
		Username:    user.Username,
		Password:    user.Password,
	}
	return data, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
