// Copyright 2016 Mirantis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"log"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/Mirantis/k8s-AppController/pkg/client"
	"github.com/Mirantis/k8s-AppController/pkg/interfaces"
	"github.com/Mirantis/k8s-AppController/pkg/report"
)

type ServiceAccount struct {
	Base
	ServiceAccount *v1.ServiceAccount
	Client         corev1.ServiceAccountInterface
}

type ExistingServiceAccount struct {
	Base
	Name   string
	Client corev1.ServiceAccountInterface
}

func serviceAccountKey(name string) string {
	return "serviceaccount/" + name
}

func (c ServiceAccount) Key() string {
	return serviceAccountKey(c.ServiceAccount.Name)
}

func serviceAccountStatus(c corev1.ServiceAccountInterface, name string) (string, error) {
	_, err := c.Get(name)
	if err != nil {
		return "error", err
	}

	return "ready", nil
}

func (c ServiceAccount) Status(meta map[string]string) (string, error) {
	return serviceAccountStatus(c.Client, c.ServiceAccount.Name)
}

func (c ServiceAccount) Create() error {
	if err := checkExistence(c); err != nil {
		log.Println("Creating ", c.Key())
		c.ServiceAccount, err = c.Client.Create(c.ServiceAccount)
		return err
	}
	return nil
}

func (c ServiceAccount) Delete() error {
	return c.Client.Delete(c.ServiceAccount.Name, &v1.DeleteOptions{})
}

func (c ServiceAccount) NameMatches(def client.ResourceDefinition, name string) bool {
	return def.ServiceAccount != nil && def.ServiceAccount.Name == name
}

func NewServiceAccount(c *v1.ServiceAccount, client corev1.ServiceAccountInterface, meta map[string]interface{}) interfaces.Resource {
	return report.SimpleReporter{BaseResource: ServiceAccount{Base: Base{meta}, ServiceAccount: c, Client: client}}
}

func NewExistingServiceAccount(name string, client corev1.ServiceAccountInterface) interfaces.Resource {
	return report.SimpleReporter{BaseResource: ExistingServiceAccount{Name: name, Client: client}}
}

// New returns a new object wrapped as Resource
func (c ServiceAccount) New(def client.ResourceDefinition, ci client.Interface) interfaces.Resource {
	return NewServiceAccount(def.ServiceAccount, ci.ServiceAccounts(), def.Meta)
}

// NewExisting returns a new object based on existing one wrapped as Resource
func (c ServiceAccount) NewExisting(name string, ci client.Interface) interfaces.Resource {
	return NewExistingServiceAccount(name, ci.ServiceAccounts())
}

func (c ExistingServiceAccount) Key() string {
	return serviceAccountKey(c.Name)
}

func (c ExistingServiceAccount) Status(meta map[string]string) (string, error) {
	return serviceAccountStatus(c.Client, c.Name)
}

func (c ExistingServiceAccount) Create() error {
	return createExistingResource(c)
}

func (c ExistingServiceAccount) Delete() error {
	return c.Client.Delete(c.Name, nil)
}
