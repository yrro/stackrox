package imageintegrations

import (
	"context"
	"fmt"

	clusterDatastore "github.com/stackrox/rox/central/cluster/datastore"
	"github.com/stackrox/rox/central/imageintegration"
	"github.com/stackrox/rox/central/imageintegration/datastore"
	countMetrics "github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/reprocessor"
	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/central/sensor/service/pipeline/reconciliation"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/integration"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/tlscheck"
	"github.com/stackrox/rox/pkg/urlfmt"
)

var (
	log = logging.LoggerForModule()
)

// Template design pattern. We define control flow here and defer logic to subclasses.
//////////////////////////////////////////////////////////////////////////////////////

// GetPipeline returns an instantiation of this particular pipeline
func GetPipeline() pipeline.Fragment {
	return NewPipeline(imageintegration.ToNotify(),
		datastore.Singleton(),
		clusterDatastore.Singleton(),
		reprocessor.Singleton())
}

// NewPipeline returns a new instance of Pipeline.
func NewPipeline(toNotify integration.ToNotify,
	datastore datastore.DataStore,
	clusterDatastore clusterDatastore.DataStore,
	enrichAndDetectLoop reprocessor.Loop) pipeline.Fragment {
	return &pipelineImpl{
		toNotify:            toNotify,
		datastore:           datastore,
		clusterDatastore:    clusterDatastore,
		enrichAndDetectLoop: enrichAndDetectLoop,
	}
}

type pipelineImpl struct {
	toNotify integration.ToNotify

	datastore           datastore.DataStore
	clusterDatastore    clusterDatastore.DataStore
	enrichAndDetectLoop reprocessor.Loop
}

func (s *pipelineImpl) Reconcile(_ context.Context, _ string, _ *reconciliation.StoreMap) error {
	// Nothing to reconcile for image integrations
	return nil
}

func (s *pipelineImpl) Match(msg *central.MsgFromSensor) bool {
	return msg.GetEvent().GetImageIntegration() != nil
}

func compareRegistries(ni, ei *storage.ImageIntegration) bool {
	return urlfmt.TrimHTTPPrefixes(ni.GetDocker().GetEndpoint()) == urlfmt.TrimHTTPPrefixes(ei.GetDocker().GetEndpoint())
}

func matchesAuth(ni, ei *storage.ImageIntegration) bool {
	return ni.GetDocker().GetUsername() == ei.GetDocker().GetUsername() &&
		ni.GetDocker().GetPassword() == ei.GetDocker().GetPassword()
}

// getMatchingImageIntegration returns the image integration that exists and should be updated
// the second return value
func (s *pipelineImpl) getMatchingImageIntegration(auto *storage.ImageIntegration, existingIntegrations []*storage.ImageIntegration) (*storage.ImageIntegration, bool) {
	var integrationToUpdate *storage.ImageIntegration
	for _, existing := range existingIntegrations {
		if !existing.GetAutogenerated() || auto.GetClusterId() != existing.GetClusterId() || !compareRegistries(auto, existing) {
			continue
		}

		// At this point, we just want to see if we already have an exact match
		// if so then we don't want to reprocess everything for no change
		if matchesAuth(auto, existing) {
			return nil, false
		}
		integrationToUpdate = existing
		break
	}
	return integrationToUpdate, true
}

// Run runs the pipeline template on the input and returns the output.
func (s *pipelineImpl) Run(ctx context.Context, clusterID string, msg *central.MsgFromSensor, _ common.MessageInjector) error {
	defer countMetrics.IncrementResourceProcessedCounter(pipeline.ActionToOperation(msg.GetEvent().GetAction()), metrics.ImageIntegration)

	clusterName, exists, err := s.clusterDatastore.GetClusterName(ctx, clusterID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cluster with id %q does not exist", clusterID)
	}

	imageIntegration := msg.GetEvent().GetImageIntegration()
	imageIntegration.ClusterId = clusterID

	validTLS, err := tlscheck.CheckTLS(imageIntegration.GetDocker().GetEndpoint())
	if err != nil {
		return err
	}

	if imageIntegration.GetDocker() == nil {
		return nil
	}

	// Using GetDocker() because the config is within a oneof
	imageIntegration.GetDocker().Insecure = !validTLS

	// Action is currently always update
	// We should not overwrite image integrations that already have a username and password
	// However, if they do not have a username and password, then we can add one that has a username and password
	existingIntegrations, err := s.datastore.GetImageIntegrations(ctx, &v1.GetImageIntegrationsRequest{})
	if err != nil {
		return err
	}

	integrationToUpdate, shouldInsert := s.getMatchingImageIntegration(imageIntegration, existingIntegrations)
	if !shouldInsert {
		return nil
	}
	if integrationToUpdate == nil {
		imageIntegration.Name = fmt.Sprintf("Autogenerated %s for cluster %s", imageIntegration.GetDocker().GetEndpoint(), clusterName)
		if err := s.toNotify.NotifyUpdated(imageIntegration); err != nil {
			return err
		}
		if _, err := s.datastore.AddImageIntegration(ctx, imageIntegration); err != nil {
			return err
		}
		// Only when adding the integration the first time do we need to run processing
		// Central receives many updates from OpenShift about the image integrations due to service accounts
		// So we can assume the other creds were valid up to this point. Also, they will eventually be picked up within an hour
		go s.enrichAndDetectLoop.ShortCircuit()
	} else {
		imageIntegration.Id = integrationToUpdate.GetId()
		imageIntegration.Name = integrationToUpdate.GetName()
		if err := s.toNotify.NotifyUpdated(imageIntegration); err != nil {
			return err
		}
		if err := s.datastore.UpdateImageIntegration(ctx, imageIntegration); err != nil {
			return err
		}
	}
	return nil
}

func (s *pipelineImpl) OnFinish(_ string) {}
