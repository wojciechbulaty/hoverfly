package modes

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/templating"
)

type HoverflySimulate interface {
	GetResponse(models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
}

type SimulateMode struct {
	Hoverfly HoverflySimulate
	MatchingStrategy string
}

func (this *SimulateMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Simulate,
		Arguments: v2.ModeArgumentsView{
			MatchingStrategy: &this.MatchingStrategy,
		},
	}
}

func (this *SimulateMode) SetArguments(arguments ModeArguments) {
	if arguments.MatchingStrategy == nil {
		this.MatchingStrategy = "strongest"
	} else {
		this.MatchingStrategy = *arguments.MatchingStrategy
	}
}

func (this SimulateMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	pair := models.RequestResponsePair{
		Request: details,
	}

	response, matchingErr := this.Hoverfly.GetResponse(details)

	if matchingErr != nil {
		return ReturnErrorAndLog(request, matchingErr, &pair, "There was an error when matching", Simulate)
	}

	response, err := templating.ApplyTemplate(request, *response) // By reference we will be mutating it

	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when templating", Simulate)
	}

	pair.Response = *response

	pair, err = this.Hoverfly.ApplyMiddleware(pair)

	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when executing middleware", Simulate)
	}

	return ReconstructResponse(request, pair), nil
}
