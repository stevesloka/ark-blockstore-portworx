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

type createSnap struct {
	Id      string `json:"id"`
	Locator `json:"locator"`
}

type createVolume struct {
	Locator `json:"locator"`
	Spec    `json:"spec"`
}

type Spec struct {
	Ephemeral        bool   `json:"ephemeral"`
	Size             int    `json:"size"`
	Format           int    `json:"format"`
	Blocksize        int    `json:"block_size"`
	HaLevel          int    `json:"ha_level"`
	Cos              int    `json:"cos"`
	IOPriority       string `json:"io_priority`
	Dedupe           bool   `json:"dedupe"`
	SnapshotInterval int    `json:"snapshot_interval"`
	Shared           bool   `json:"shared"`
}

type Locator struct {
	Name string `json:"name"`
}
