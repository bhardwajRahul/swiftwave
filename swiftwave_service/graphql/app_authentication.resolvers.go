package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	haproxymanager "github.com/swiftwave-org/swiftwave/haproxy_manager"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
	"gorm.io/gorm"
)

// Users is the resolver for the users field.
func (r *appBasicAuthAccessControlListResolver) Users(ctx context.Context, obj *model.AppBasicAuthAccessControlList) ([]*model.AppBasicAuthAccessControlUser, error) {
	records, err := core.FetchAppBasicAuthAccessControlUsers(ctx, &r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	users := make([]*model.AppBasicAuthAccessControlUser, len(records))
	for i, record := range records {
		users[i] = appBasicAuthAccessControlUserToGraphqlObject(record)
	}
	return users, nil
}

// CreateAppBasicAuthAccessControlList is the resolver for the createAppBasicAuthAccessControlList field.
func (r *mutationResolver) CreateAppBasicAuthAccessControlList(ctx context.Context, input model.AppBasicAuthAccessControlListInput) (*model.AppBasicAuthAccessControlList, error) {
	tx := r.ServiceManager.DbClient.Begin()
	defer tx.Rollback()

	// add in db
	appBasicAuthAccessControlList := appBasicAuthAccessControlListInputToDatabaseObject(&input)
	err := appBasicAuthAccessControlList.Create(ctx, tx)
	if err != nil {
		return nil, err
	}
	// add in haproxy + commit
	ctx = context.WithValue(ctx, "access_control_user_list_name", appBasicAuthAccessControlList.GeneratedName)
	err = r.RunActionsInAllHAProxyNodes(ctx, tx, func(ctx context.Context, db *gorm.DB, transactionId string, manager *haproxymanager.Manager) error {
		userList := ctx.Value("access_control_user_list_name").(string)
		return manager.AddUserList(transactionId, userList)
	})
	if err != nil {
		return nil, err
	}
	// commit to db
	return appBasicAuthAccessControlListToGraphqlObject(appBasicAuthAccessControlList), tx.Commit().Error
}

// DeleteAppBasicAuthAccessControlList is the resolver for the deleteAppBasicAuthAccessControlList field.
func (r *mutationResolver) DeleteAppBasicAuthAccessControlList(ctx context.Context, id uint) (bool, error) {
	tx := r.ServiceManager.DbClient.Begin()
	defer tx.Rollback()

	// delete in db
	record := &core.AppBasicAuthAccessControlList{}
	err := record.FindById(ctx, tx, id)
	if err != nil {
		return false, err
	}

	err = record.Delete(ctx, tx)
	if err != nil {
		return false, err
	}

	// delete in haproxy + commit
	ctx = context.WithValue(ctx, "access_control_user_list_name", record.GeneratedName)
	err = r.RunActionsInAllHAProxyNodes(ctx, tx, func(ctx context.Context, db *gorm.DB, transactionId string, manager *haproxymanager.Manager) error {
		userList := ctx.Value("access_control_user_list_name").(string)
		return manager.DeleteUserList(transactionId, userList)
	})
	if err != nil {
		return false, err
	}

	// commit to db
	return true, tx.Commit().Error
}

// CreateAppBasicAuthAccessControlUser is the resolver for the createAppBasicAuthAccessControlUser field.
func (r *mutationResolver) CreateAppBasicAuthAccessControlUser(ctx context.Context, input model.AppBasicAuthAccessControlUserInput) (*model.AppBasicAuthAccessControlUser, error) {
	tx := r.ServiceManager.DbClient.Begin()
	defer tx.Rollback()

	// find in db
	userList := &core.AppBasicAuthAccessControlList{}
	err := userList.FindById(ctx, tx, input.AppBasicAuthAccessControlListID)
	if err != nil {
		return nil, err
	}

	// add in db
	record := appBasicAuthAccessControlUserInputToDatabaseObject(&input)
	encryptedPassword, err := haproxymanager.GenerateSecuredPasswordForBasicAuthentication(record.PlainTextPassword)
	if err != nil {
		return nil, err
	}
	record.EncryptedPassword = encryptedPassword
	err = record.Create(ctx, tx)
	if err != nil {
		return nil, err
	}

	// add in haproxy + commit
	ctx = context.WithValue(ctx, "access_control_user_list_name", userList.GeneratedName)
	ctx = context.WithValue(ctx, "access_control_user_name", record.Username)
	ctx = context.WithValue(ctx, "access_control_user_password", record.EncryptedPassword)
	err = r.RunActionsInAllHAProxyNodes(ctx, tx, func(ctx context.Context, db *gorm.DB, transactionId string, manager *haproxymanager.Manager) error {
		userListName := ctx.Value("access_control_user_list_name").(string)
		userName := ctx.Value("access_control_user_name").(string)
		password := ctx.Value("access_control_user_password").(string)
		return manager.AddUserInUserList(transactionId, userListName, userName, password)
	})
	if err != nil {
		return nil, err
	}

	return appBasicAuthAccessControlUserToGraphqlObject(record), tx.Commit().Error
}

// UpdateAppBasicAuthAccessControlUserPassword is the resolver for the updateAppBasicAuthAccessControlUserPassword field.
func (r *mutationResolver) UpdateAppBasicAuthAccessControlUserPassword(ctx context.Context, id uint, password string) (bool, error) {
	tx := r.ServiceManager.DbClient.Begin()
	defer tx.Rollback()

	// find in db
	record := &core.AppBasicAuthAccessControlUser{}
	err := record.FindById(ctx, tx, id)
	if err != nil {
		return false, err
	}

	// update in db
	record.PlainTextPassword = password
	encryptedPassword, err := haproxymanager.GenerateSecuredPasswordForBasicAuthentication(password)
	if err != nil {
		return false, err
	}
	record.EncryptedPassword = encryptedPassword
	err = record.Update(ctx, tx)
	if err != nil {
		return false, err
	}

	// find user list in db
	userList := &core.AppBasicAuthAccessControlList{}
	err = userList.FindById(ctx, tx, record.AppBasicAuthAccessControlListID)
	if err != nil {
		return false, err
	}

	// update in haproxy + commit
	ctx = context.WithValue(ctx, "access_control_user_list_name", userList.GeneratedName)
	ctx = context.WithValue(ctx, "access_control_user_name", record.Username)
	ctx = context.WithValue(ctx, "access_control_user_password", record.EncryptedPassword)
	err = r.RunActionsInAllHAProxyNodes(ctx, tx, func(ctx context.Context, db *gorm.DB, transactionId string, manager *haproxymanager.Manager) error {
		userListName := ctx.Value("access_control_user_list_name").(string)
		userName := ctx.Value("access_control_user_name").(string)
		password := ctx.Value("access_control_user_password").(string)
		return manager.ChangeUserPasswordInUserList(transactionId, userListName, userName, password)
	})
	if err != nil {
		return false, err
	}

	err = tx.Commit().Error
	return err == nil, err
}

// DeleteAppBasicAuthAccessControlUser is the resolver for the deleteAppBasicAuthAccessControlUser field.
func (r *mutationResolver) DeleteAppBasicAuthAccessControlUser(ctx context.Context, id uint) (bool, error) {
	tx := r.ServiceManager.DbClient.Begin()
	defer tx.Rollback()

	record := &core.AppBasicAuthAccessControlUser{}
	err := record.FindById(ctx, tx, id)
	if err != nil {
		return false, err
	}

	userList := &core.AppBasicAuthAccessControlList{}
	err = userList.FindById(ctx, tx, record.AppBasicAuthAccessControlListID)
	if err != nil {
		return false, err
	}

	// delete from db
	err = record.Delete(ctx, tx)
	if err != nil {
		return false, err
	}

	// delete from haproxy + commit
	ctx = context.WithValue(ctx, "access_control_user_list_name", userList.GeneratedName)
	ctx = context.WithValue(ctx, "access_control_user_name", record.Username)
	err = r.RunActionsInAllHAProxyNodes(ctx, tx, func(ctx context.Context, db *gorm.DB, transactionId string, manager *haproxymanager.Manager) error {
		userListName := ctx.Value("access_control_user_list_name").(string)
		userName := ctx.Value("access_control_user_name").(string)
		return manager.DeleteUserFromUserList(transactionId, userListName, userName)
	})
	if err != nil {
		return false, err
	}

	err = tx.Commit().Error
	return err == nil, err
}

// AppBasicAuthAccessControlLists is the resolver for the appBasicAuthAccessControlLists field.
func (r *queryResolver) AppBasicAuthAccessControlLists(ctx context.Context) ([]*model.AppBasicAuthAccessControlList, error) {
	records, err := core.FindAllAppBasicAuthAccessControlLists(ctx, &r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	appBasicAuthAccessControlLists := make([]*model.AppBasicAuthAccessControlList, len(records))
	for i, record := range records {
		appBasicAuthAccessControlLists[i] = appBasicAuthAccessControlListToGraphqlObject(record)
	}
	return appBasicAuthAccessControlLists, nil
}

// AppBasicAuthAccessControlList returns AppBasicAuthAccessControlListResolver implementation.
func (r *Resolver) AppBasicAuthAccessControlList() AppBasicAuthAccessControlListResolver {
	return &appBasicAuthAccessControlListResolver{r}
}

type appBasicAuthAccessControlListResolver struct{ *Resolver }