package client

import (
	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cl *Client) RetrieveBrokers() ([]v1beta1.ClusterServiceBroker, error) {
	brokers, err := cl.ServicecatalogV1beta1().ClusterServiceBrokers().List(v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list brokers (%s)", err)
	}

	return brokers.Items, nil
}
func (cl *Client) RetrieveBroker(name string) (*v1beta1.ClusterServiceBroker, error) {
	broker, err := cl.ServicecatalogV1beta1().ClusterServiceBrokers().Get(name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get broker '%s' (%s)", name, err)
	}

	return broker, nil
}

// RetrieveBrokerByClass retrieves the parent broker of a class.
func (cl *Client) RetrieveBrokerByClass(class *v1beta1.ClusterServiceClass,
) (*v1beta1.ClusterServiceBroker, error) {
	brokerName := class.Spec.ClusterServiceBrokerName
	broker, err := cl.ServicecatalogV1beta1().ClusterServiceBrokers().Get(brokerName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return broker, nil
}

func (cl *Client) Sync(name string, retries int) error {
	for j := 0; j < retries; j++ {
		catalog, err := cl.RetrieveBroker(name)
		if err != nil {
			return err
		}

		catalog.Spec.RelistRequests = catalog.Spec.RelistRequests + 1

		_, err = cl.ServicecatalogV1beta1().ClusterServiceBrokers().Update(catalog)
		if err == nil {
			return nil
		}
		if !errors.IsConflict(err) {
			return fmt.Errorf("could not sync service broker (%s)", err)
		}
	}

	return fmt.Errorf("could not sync service broker after %d tries", retries)
}
