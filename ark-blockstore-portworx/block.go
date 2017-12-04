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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/heptio/ark/pkg/cloudprovider"
	"github.com/heptio/ark/pkg/util/collections"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	apiURL       = "http://portworx-service.kube-system:9001"
	volumeIDPath = "spec.portworxVolume.volumeID"
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
	f.log.Info("BlockStore.Init called")
	return nil
}

// CreateVolumeFromSnapshot creates a new block volume in the specified
// availability zone, initialized from the provided snapshot,
// and with the specified type and IOPS (if using provisioned IOPS).
func (f *BlockStore) CreateVolumeFromSnapshot(snapshotID, volumeType, volumeAZ string, iops *int64) (string, error) {

	f.log.Info("CreateVolumeFromSnapshot: snapshotID: ", snapshotID)

	body := createVolume{
		Locator: Locator{
			Name: snapshotID,
		},
		Spec: Spec{
			Ephemeral: false,
			Format:    2,
			Size:      2147483648,
		},
	}

	httpURL := fmt.Sprintf("%s/v1/osd-volumes", apiURL)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)

	res, err := http.Post(httpURL, "application/json", b)

	if err != nil {
		f.log.Error("Error creating volume from snapshot! ", err)
		return "", err
	}

	f.log.Info("response: ", res.StatusCode)

	var data map[string]interface{}
	bodyResponse, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyResponse, &data)

	if err != nil {
		f.log.Error("Error unmarshaling json ", err)
		return "", err
	}

	volID := data["id"].(string)

	return volID, nil
}

// GetVolumeID returns the cloud provider specific identifier for the PersistentVolume.
func (f *BlockStore) GetVolumeID(pv runtime.Unstructured) (string, error) {

	f.log.Info("GetVolumeID")

	if !collections.Exists(pv.UnstructuredContent(), volumeIDPath) {
		return "", nil
	}

	volumeID, err := collections.GetString(pv.UnstructuredContent(), volumeIDPath)
	if err != nil {
		return "", err
	}

	f.log.Info("GetVolumeID:", volumeID)

	return volumeID, nil
}

// SetVolumeID sets the cloud provider specific identifier for the PersistentVolume.
func (f *BlockStore) SetVolumeID(pv runtime.Unstructured, volumeID string) (runtime.Unstructured, error) {

	px, err := collections.GetMap(pv.UnstructuredContent(), volumeIDPath)
	if err != nil {
		return nil, err
	}

	px["volumeID"] = volumeID

	return pv, nil
}

// GetVolumeInfo returns the type and IOPS (if using provisioned IOPS) for
// the specified block volume in the given availability zone.
func (f *BlockStore) GetVolumeInfo(volumeID, volumeAZ string) (string, *int64, error) {
	return "portworx", nil, nil
}

// IsVolumeReady returns whether the specified volume is ready to be used.
func (f *BlockStore) IsVolumeReady(volumeID, volumeAZ string) (bool, error) {

	data, err := f.getVolumeInfo(volumeID)

	if err != nil {
		return false, err
	}

	if data["status"] != "up" {
		return false, nil
	}

	return true, nil
}

// CreateSnapshot creates a snapshot of the specified block volume, and applies the provided
// set of tags to the snapshot.
func (f *BlockStore) CreateSnapshot(volumeID, volumeAZ string, tags map[string]string) (string, error) {

	f.log.Info("CreateSnapshot")

	snapshotID := uuid.NewV4().String()
	body := createSnap{
		Id: volumeID,
		Locator: Locator{
			Name: snapshotID,
		},
	}

	httpURL := fmt.Sprintf("%s/v1/osd-snapshot", apiURL)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)

	_, err := http.Post(httpURL, "application/json", b)

	if err != nil {
		f.log.Error("Error creating snapshot: ", err)
		return "", err
	}

	f.log.Info("snapshotID: ", snapshotID)

	return snapshotID, nil
}

// DeleteSnapshot deletes the specified volume snapshot.
func (f *BlockStore) DeleteSnapshot(snapshotID string) error {

	httpURL := fmt.Sprintf("%s/v1/osd-volumes/%s", apiURL, snapshotID)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", httpURL, nil)
	resp, err := client.Do(req)

	// Process response
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
