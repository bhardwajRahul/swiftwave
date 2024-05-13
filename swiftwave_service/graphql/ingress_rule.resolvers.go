package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.46

import (
	"context"
	"errors"
	"strings"

	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
)

// Domain is the resolver for the domain field.
func (r *ingressRuleResolver) Domain(ctx context.Context, obj *model.IngressRule) (*model.Domain, error) {
	domain := &core.Domain{}
	if obj.DomainID == nil {
		return nil, nil
	}
	err := domain.FindById(ctx, r.ServiceManager.DbClient, *obj.DomainID)
	if err != nil {
		return nil, err
	}
	return domainToGraphqlObject(domain), nil
}

// Application is the resolver for the application field.
func (r *ingressRuleResolver) Application(ctx context.Context, obj *model.IngressRule) (*model.Application, error) {
	application := &core.Application{}
	if obj.TargetType == model.IngressRuleTargetTypeExternalService {
		return applicationToGraphqlObject(application), nil
	}
	err := application.FindById(ctx, r.ServiceManager.DbClient, obj.ApplicationID)
	if err != nil {
		return nil, err
	}
	return applicationToGraphqlObject(application), nil
}

// CreateIngressRule is the resolver for the createIngressRule field.
func (r *mutationResolver) CreateIngressRule(ctx context.Context, input model.IngressRuleInput) (*model.IngressRule, error) {
	record := ingressRuleInputToDatabaseObject(&input)
	if record.TargetType == core.ExternalServiceIngressRule && strings.Compare(record.ExternalService, "") == 0 {
		return nil, errors.New("external service is required")
	}
	restrictedPorts := make([]int, 0)
	for _, port := range r.Config.SystemConfig.RestrictedPorts {
		restrictedPorts = append(restrictedPorts, int(port))
	}
	err := record.Create(ctx, r.ServiceManager.DbClient, restrictedPorts)
	if err != nil {
		return nil, err
	}
	// schedule task
	err = r.WorkerManager.EnqueueIngressRuleApplyRequest(record.ID)
	if err != nil {
		return nil, errors.New("failed to schedule task to apply ingress rule")
	}
	return ingressRuleToGraphqlObject(record), nil
}

// DeleteIngressRule is the resolver for the deleteIngressRule field.
func (r *mutationResolver) DeleteIngressRule(ctx context.Context, id uint) (bool, error) {
	record := core.IngressRule{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	err = record.Delete(ctx, r.ServiceManager.DbClient, false)
	if err != nil {
		if errors.Is(err, core.IngressRuleDeletingError) {
			_ = r.WorkerManager.EnqueueIngressRuleDeleteRequest(record.ID)
		}
		return false, err
	}
	// schedule task
	err = r.WorkerManager.EnqueueIngressRuleDeleteRequest(record.ID)
	if err != nil {
		return false, errors.New("failed to schedule task to delete ingress rule")
	}
	return true, nil
}

// IngressRule is the resolver for the ingressRule field.
func (r *queryResolver) IngressRule(ctx context.Context, id uint) (*model.IngressRule, error) {
	record := core.IngressRule{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	return ingressRuleToGraphqlObject(&record), nil
}

// IngressRules is the resolver for the ingressRules field.
func (r *queryResolver) IngressRules(ctx context.Context) ([]*model.IngressRule, error) {
	records, err := core.FindAllIngressRules(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	var result []*model.IngressRule
	for _, record := range records {
		result = append(result, ingressRuleToGraphqlObject(record))
	}
	return result, nil
}

// IsNewIngressRuleValid is the resolver for the isNewIngressRuleValid field.
func (r *queryResolver) IsNewIngressRuleValid(ctx context.Context, input model.IngressRuleValidationInput) (bool, error) {
	record := ingressRuleValidationInputToDatabaseObject(&input)
	restrictedPorts := make([]int, 0)
	for _, port := range r.Config.SystemConfig.RestrictedPorts {
		restrictedPorts = append(restrictedPorts, int(port))
	}
	if err := record.IsValidNewIngressRule(ctx, r.ServiceManager.DbClient, restrictedPorts); err != nil {
		return false, err
	}
	return true, nil
}

// IngressRule returns IngressRuleResolver implementation.
func (r *Resolver) IngressRule() IngressRuleResolver { return &ingressRuleResolver{r} }

type ingressRuleResolver struct{ *Resolver }
