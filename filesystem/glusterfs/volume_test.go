// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package glusterfs

import (
	"fmt"
	"testing"
)

/*
func TestGetAllVolume(t *testing.T) {
	glusterfsVolumeControl, _ := CreateGlusterfsVolumeControl()
	fmt.Println(glusterfsVolumeControl.GetAllVolume())
}

func TestGetAvailableHost(t *testing.T) {
	glusterfsVolumeControl, _ := CreateGlusterfsVolumeControl()
	fmt.Println(glusterfsVolumeControl.getAvailableHost())
}
*/

func TestGetHostStatus(t *testing.T) {
	glusterfsVolumeControl, _ := CreateGlusterfsVolumeControl()
	fmt.Println(glusterfsVolumeControl.getHostStatus())
}