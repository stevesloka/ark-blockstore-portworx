/*
Copyright 2017 the Heptio Ark contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/heptio/ark/pkg/cloudprovider"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	portworxAPI = "http://portworx-service.kube-system:9001"
)

type BlockStore struct {
	log logrus.FieldLogger
}

func NewBlockStore(log logrus.FieldLogger) cloudprovider.BlockStore {
	return &BlockStore{
		log: log,
	}
}

func (f *BlockStore) Init(config map[string]string) error {
	f.log.Infof("BlockStore.Init called")
	return nil
}

// CreateVolumeFromSnapshot creates a new block volume in the specified
// availability zone, initialized from the provided snapshot,
// and with the specified type and IOPS (if using provisioned IOPS).
func (f *BlockStore) CreateVolumeFromSnapshot(snapshotID, volumeType, volumeAZ string, iops *int64) (volumeID string, err error) {
	return "", nil
}

// GetVolumeID returns the cloud provider specific identifier for the PersistentVolume.
func (f *BlockStore) GetVolumeID(pv runtime.Unstructured) (string, error) {
	return "", nil
}

// SetVolumeID sets the cloud provider specific identifier for the PersistentVolume.
func (f *BlockStore) SetVolumeID(pv runtime.Unstructured, volumeID string) (runtime.Unstructured, error) {
	return pv, nil
}

// GetVolumeInfo returns the type and IOPS (if using provisioned IOPS) for
// the specified block volume in the given availability zone.
func (f *BlockStore) GetVolumeInfo(volumeID, volumeAZ string) (string, *int64, error) {
	return "", nil, nil
}

// IsVolumeReady returns whether the specified volume is ready to be used.
func (f *BlockStore) IsVolumeReady(volumeID, volumeAZ string) (ready bool, err error) {
	return true, nil
}

// CreateSnapshot creates a snapshot of the specified block volume, and applies the provided
// set of tags to the snapshot.
func (f *BlockStore) CreateSnapshot(volumeID, volumeAZ string, tags map[string]string) (snapshotID string, err error) {
	return "", nil
}

// DeleteSnapshot deletes the specified volume snapshot.
func (f *BlockStore) DeleteSnapshot(snapshotID string) error {
	return nil
}
