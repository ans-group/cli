package pss

import (
	"github.com/ans-group/sdk-go/pkg/service/pss"
)

type RequestCollection []pss.Request

func (m RequestCollection) DefaultColumns() []string {
	return []string{"id", "author_name", "type", "subject", "status", "last_replied_at", "created_at"}
}

type ReplyCollection []pss.Reply

func (m ReplyCollection) DefaultColumns() []string {
	return []string{"id", "author_name", "created_at"}
}

type FeedbackCollection []pss.Feedback

func (m FeedbackCollection) DefaultColumns() []string {
	return []string{"id", "score", "created_at"}
}

type CaseOptionCollection []pss.CaseOption

func (m CaseOptionCollection) DefaultColumns() []string {
	return []string{"option"}
}

type IncidentCaseCollection []pss.IncidentCase

func (m IncidentCaseCollection) DefaultColumns() []string {
	return []string{"id", "title", "status", "created_at", "updated_at"}
}

type ChangeCaseCollection []pss.ChangeCase

func (m ChangeCaseCollection) DefaultColumns() []string {
	return []string{"id", "title", "status", "created_at", "updated_at"}
}

type ProblemCaseCollection []pss.ProblemCase

func (m ProblemCaseCollection) DefaultColumns() []string {
	return []string{"id", "title", "status", "created_at", "updated_at"}
}

type CaseCategoryCollection []pss.CaseCategory

func (m CaseCategoryCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

type SupportedServiceCollection []pss.SupportedService

func (m SupportedServiceCollection) DefaultColumns() []string {
	return []string{"id", "name"}
}

type CaseUpdateCollection []pss.CaseUpdate

func (m CaseUpdateCollection) DefaultColumns() []string {
	return []string{"id", "subject", "created_at", "updated_at"}
}

type CaseCollection []pss.Case

func (m CaseCollection) DefaultColumns() []string {
	return []string{"id", "case_type", "title", "status", "created_at", "updated_at"}
}
