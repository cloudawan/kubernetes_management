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

package notification

/*
import (
	"fmt"
	"testing"
	"time"
)
*/
/*
func TestSave(t *testing.T) {

	indicatorSlice := make([]Indicator, 1)
	indicatorSlice[0] = Indicator{"cpu", false, 0.3, 800000000, false, 0.3, 1000000}
	notifierSlice := make([]Notifier, 1)
	notifierSlice[0] = NotifierEmail{[]string{"cloudawanemailtest@gmail.com"}}
	replicationControllerNotifier := ReplicationControllerNotifier{
		true, 10 * time.Second, 0 * time.Second, "192.168.0.33", 8080,
		"default", "replicationController", "flask", notifierSlice, indicatorSlice}

	replicationControllerNotifierSerializable := ConvertToSerializable(replicationControllerNotifier)

	fmt.Println(SaveReplicationControllerNotifierSerializable(&replicationControllerNotifierSerializable))
}
*/
/*
func TestDelete(t *testing.T) {
	fmt.Println(DeleteReplicationControllerNotifierSerializable("default", "replicationController", "flask"))
}
*/
/*
func TestLoad(t *testing.T) {

	replicationControllerNotifierSerializable, err := LoadReplicationControllerNotifierSerializable("default", "replicationController", "flask")
	fmt.Println(replicationControllerNotifierSerializable, err)
}
*/
/*
func TestLoadAll(t *testing.T) {
	fmt.Println(LoadAllReplicationControllerNotifierSerializable())
}
*/
