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
	"bufio"
	"bytes"
	"errors"
	"github.com/cloudawan/cloudone/utility/configuration"
	"github.com/cloudawan/cloudone_utility/logger"
	"github.com/cloudawan/cloudone_utility/sshclient"
	"strconv"
	"strings"
	"time"
)

func CreateGlusterfsVolumeControl() (*GlusterfsVolumeControl, error) {
	var ok bool

	glusterfsHostSlice, ok := configuration.LocalConfiguration.GetStringSlice("glusterfsHost")
	if ok == false {
		log.Error("Can't load glusterfsClusterIPSlice")
		return nil, errors.New("Can't load glusterfsClusterIPSlice")
	}

	glusterfsPath, ok := configuration.LocalConfiguration.GetString("glusterfsPath")
	if ok == false {
		log.Error("Can't load glusterfsPath")
		return nil, errors.New("Can't load glusterfsPath")
	}

	glusterfsSSHTimeoutInSecond, ok := configuration.LocalConfiguration.GetInt("glusterfsSSHTimeoutInSecond")
	if ok == false {
		log.Error("Can't load glusterfsSSHTimeoutInSecond")
		return nil, errors.New("Can't load glusterfsSSHTimeoutInSecond")
	}

	glusterfsSSHHostSlice, ok := configuration.LocalConfiguration.GetStringSlice("glusterfsSSHHost")
	if ok == false {
		log.Error("Can't load glusterfsSSHHost")
		return nil, errors.New("Can't load glusterfsSSHHost")
	}

	glusterfsSSHPort, ok := configuration.LocalConfiguration.GetInt("glusterfsSSHPort")
	if ok == false {
		log.Error("Can't load glusterfsSSHPort")
		return nil, errors.New("Can't load glusterfsSSHPort")
	}

	glusterfsSSHUser, ok := configuration.LocalConfiguration.GetString("glusterfsSSHUser")
	if ok == false {
		log.Error("Can't load glusterfsSSHUser")
		return nil, errors.New("Can't load glusterfsSSHUser")
	}

	glusterfsSSHPassword, ok := configuration.LocalConfiguration.GetString("glusterfsSSHPassword")
	if ok == false {
		log.Error("Can't load glusterfsSSHPassword")
		return nil, errors.New("Can't load glusterfsSSHPassword")
	}

	glusterfsVolumeControl := &GlusterfsVolumeControl{
		glusterfsHostSlice,
		glusterfsPath,
		glusterfsSSHTimeoutInSecond,
		glusterfsSSHHostSlice,
		glusterfsSSHPort,
		glusterfsSSHUser,
		glusterfsSSHPassword,
	}

	return glusterfsVolumeControl, nil
}

type GlusterfsVolume struct {
	VolumeName     string
	Type           string
	VolumeID       string
	Status         string
	NumberOfBricks string
	TransportType  string
	Bricks         []string
	Size           int
}

type GlusterfsVolumeControl struct {
	GlusterfsClusterHostSlice   []string
	GlusterfsPath               string
	GlusterfsSSHTimeoutInSecond int
	GlusterfsSSHHostSlice       []string
	GlusterfsSSHPort            int
	GlusterfsSSHUser            string
	GlusterfsSSHPassword        string
}

func parseVolumeInfo(text string) (returnedGlusterfsVolumeSlice []GlusterfsVolume, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("parseVolumeInfo Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedGlusterfsVolumeSlice = nil
			returnedError = err.(error)
		}
	}()

	glusterfsVolumeSlice := make([]GlusterfsVolume, 0)

	scanner := bufio.NewScanner(bytes.NewBufferString(text))
	var glusterfsVolume *GlusterfsVolume = nil
	for scanner.Scan() {
		line := scanner.Text()
		if line == " " {
			if glusterfsVolume != nil {
				glusterfsVolumeSlice = append(glusterfsVolumeSlice, *glusterfsVolume)
			}
			glusterfsVolume = &GlusterfsVolume{}
		} else if strings.HasPrefix(line, "Volume Name: ") {
			glusterfsVolume.VolumeName = line[len("Volume Name: "):]
		} else if strings.HasPrefix(line, "Type: ") {
			glusterfsVolume.Type = line[len("Type: "):]
		} else if strings.HasPrefix(line, "Volume ID: ") {
			glusterfsVolume.VolumeID = line[len("Volume ID: "):]
		} else if strings.HasPrefix(line, "Status: ") {
			glusterfsVolume.Status = line[len("Status: "):]
		} else if strings.HasPrefix(line, "Number of Bricks: ") {
			glusterfsVolume.NumberOfBricks = line[len("Number of Bricks: "):]
			var err error
			glusterfsVolume.Size, err = strconv.Atoi(strings.TrimSpace(strings.Split(line, "=")[1]))
			if err != nil {
				log.Error("Parse brick error %s", err)
				return nil, err
			}
		} else if strings.HasPrefix(line, "Transport-type: ") {
			glusterfsVolume.TransportType = line[len("Transport-type: "):]
		} else if line == "Bricks:" {
			brickSlice := make([]string, 0)
			for i := 0; i < glusterfsVolume.Size; i++ {
				scanner.Scan()
				brickSlice = append(brickSlice, scanner.Text())
			}
			glusterfsVolume.Bricks = brickSlice
		} else {
			// Ignore unexpected data
		}
	}
	if glusterfsVolume != nil {
		glusterfsVolumeSlice = append(glusterfsVolumeSlice, *glusterfsVolume)
	}
	if err := scanner.Err(); err != nil {
		log.Error("Scanner error %s", err)
		return nil, err
	}
	return glusterfsVolumeSlice, nil
}

func (glusterfsVolumeControl *GlusterfsVolumeControl) getAvailableHost() (*string, error) {
	for _, host := range glusterfsVolumeControl.GlusterfsClusterHostSlice {
		_, err := sshclient.InteractiveSSH(
			time.Duration(glusterfsVolumeControl.GlusterfsSSHTimeoutInSecond)*time.Second,
			host,
			glusterfsVolumeControl.GlusterfsSSHPort,
			glusterfsVolumeControl.GlusterfsSSHUser,
			glusterfsVolumeControl.GlusterfsSSHPassword,
			nil,
			nil)
		if err == nil {
			return &host, nil
		}
	}

	return nil, errors.New("No available host")
}

func (glusterfsVolumeControl *GlusterfsVolumeControl) GetAllVolume() ([]GlusterfsVolume, error) {
	host, err := glusterfsVolumeControl.getAvailableHost()
	if err != nil {
		return nil, err
	}

	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, "sudo gluster volume info\n")

	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = "cloud4win\n"

	resultSlice, err := sshclient.InteractiveSSH(
		time.Duration(glusterfsVolumeControl.GlusterfsSSHTimeoutInSecond)*time.Second,
		*host,
		glusterfsVolumeControl.GlusterfsSSHPort,
		glusterfsVolumeControl.GlusterfsSSHUser,
		glusterfsVolumeControl.GlusterfsSSHPassword,
		commandSlice,
		interactiveMap)

	glusterfsVolumeSlice, err := parseVolumeInfo(resultSlice[0])
	if err != nil {
		log.Error("Parse volume info error %s", err)
		return nil, err
	} else {
		return glusterfsVolumeSlice, nil
	}
}

func (glusterfsVolumeControl *GlusterfsVolumeControl) CreateVolume(name string,
	stripe int, replica int, transport string, hostSlice []string) error {
	if stripe == 0 {
		if replica != len(hostSlice) {
			return errors.New("Replica amount is not the same as ip amount")
		}
	} else if replica == 0 {
		if stripe != len(hostSlice) {
			return errors.New("Stripe amount is not the same as ip amount")
		}
	} else {
		if stripe*replica != len(hostSlice) {
			return errors.New("Replica * Stripe amount is not the same as ip amount")
		}
	}

	host, err := glusterfsVolumeControl.getAvailableHost()
	if err != nil {
		return err
	}

	commandBuffer := bytes.Buffer{}
	commandBuffer.WriteString("sudo gluster --mode=script volume create ")
	commandBuffer.WriteString(name)

	if stripe > 0 {
		commandBuffer.WriteString(" stripe ")
		commandBuffer.WriteString(strconv.Itoa(stripe))
	}
	if replica > 0 {
		commandBuffer.WriteString(" replica ")
		commandBuffer.WriteString(strconv.Itoa(replica))
	}
	commandBuffer.WriteString(" transport ")
	commandBuffer.WriteString(transport)
	for _, ip := range hostSlice {
		path := " " + ip + ":" + glusterfsVolumeControl.GlusterfsPath + "/" + name
		commandBuffer.WriteString(path)
	}
	commandBuffer.WriteString(" force\n")
	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, commandBuffer.String())

	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = glusterfsVolumeControl.GlusterfsSSHPassword + "\n"

	resultSlice, err := sshclient.InteractiveSSH(
		time.Duration(glusterfsVolumeControl.GlusterfsSSHTimeoutInSecond)*time.Second,
		*host,
		glusterfsVolumeControl.GlusterfsSSHPort,
		glusterfsVolumeControl.GlusterfsSSHUser,
		glusterfsVolumeControl.GlusterfsSSHPassword,
		commandSlice,
		interactiveMap)

	if err != nil {
		log.Error("Create volume error %s resultSlice %s", err, resultSlice)
		return err
	} else {
		if strings.Contains(resultSlice[0], "success") {
			return nil
		} else {
			log.Debug("Issue command: " + commandBuffer.String())
			log.Error("Fail to create volume with error: " + resultSlice[0])
			return errors.New(resultSlice[0])
		}
	}
}

func (glusterfsVolumeControl *GlusterfsVolumeControl) StartVolume(name string) error {
	host, err := glusterfsVolumeControl.getAvailableHost()
	if err != nil {
		return err
	}

	commandBuffer := bytes.Buffer{}
	commandBuffer.WriteString("sudo gluster --mode=script volume start ")
	commandBuffer.WriteString(name)
	commandBuffer.WriteString(" \n")
	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, commandBuffer.String())

	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = glusterfsVolumeControl.GlusterfsSSHPassword + "\n"

	resultSlice, err := sshclient.InteractiveSSH(
		time.Duration(glusterfsVolumeControl.GlusterfsSSHTimeoutInSecond)*time.Second,
		*host,
		glusterfsVolumeControl.GlusterfsSSHPort,
		glusterfsVolumeControl.GlusterfsSSHUser,
		glusterfsVolumeControl.GlusterfsSSHPassword,
		commandSlice,
		interactiveMap)

	if err != nil {
		log.Error("Create volume error %s resultSlice %s", err, resultSlice)
		return err
	} else {
		if strings.Contains(resultSlice[0], "success") {
			return nil
		} else {
			log.Debug("Issue command: " + commandBuffer.String())
			log.Error("Fail to start volume with error: " + resultSlice[0])
			return errors.New(resultSlice[0])
		}
	}
}

func (glusterfsVolumeControl *GlusterfsVolumeControl) StopVolume(name string) error {
	host, err := glusterfsVolumeControl.getAvailableHost()
	if err != nil {
		return err
	}

	commandBuffer := bytes.Buffer{}
	commandBuffer.WriteString("sudo gluster --mode=script volume stop ")
	commandBuffer.WriteString(name)
	commandBuffer.WriteString(" \n")
	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, commandBuffer.String())

	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = glusterfsVolumeControl.GlusterfsSSHPassword + "\n"

	resultSlice, err := sshclient.InteractiveSSH(
		time.Duration(glusterfsVolumeControl.GlusterfsSSHTimeoutInSecond)*time.Second,
		*host,
		glusterfsVolumeControl.GlusterfsSSHPort,
		glusterfsVolumeControl.GlusterfsSSHUser,
		glusterfsVolumeControl.GlusterfsSSHPassword,
		commandSlice,
		interactiveMap)

	if err != nil {
		log.Error("Create volume error %s resultSlice %s", err, resultSlice)
		return err
	} else {
		if strings.Contains(resultSlice[0], "success") {
			return nil
		} else {
			log.Debug("Issue command: " + commandBuffer.String())
			log.Error("Fail to stop volume with error: " + resultSlice[0])
			return errors.New(resultSlice[0])
		}
	}
}

func (glusterfsVolumeControl *GlusterfsVolumeControl) DeleteVolume(name string) error {
	host, err := glusterfsVolumeControl.getAvailableHost()
	if err != nil {
		return err
	}

	commandBuffer := bytes.Buffer{}
	commandBuffer.WriteString("sudo gluster --mode=script volume delete ")
	commandBuffer.WriteString(name)
	commandBuffer.WriteString(" \n")
	commandSlice := make([]string, 0)
	commandSlice = append(commandSlice, commandBuffer.String())

	interactiveMap := make(map[string]string)
	interactiveMap["[sudo]"] = glusterfsVolumeControl.GlusterfsSSHPassword + "\n"

	resultSlice, err := sshclient.InteractiveSSH(
		time.Duration(glusterfsVolumeControl.GlusterfsSSHTimeoutInSecond)*time.Second,
		*host,
		glusterfsVolumeControl.GlusterfsSSHPort,
		glusterfsVolumeControl.GlusterfsSSHUser,
		glusterfsVolumeControl.GlusterfsSSHPassword,
		commandSlice,
		interactiveMap)

	if err != nil {
		log.Error("Create volume error %s resultSlice %s", err, resultSlice)
		return err
	} else {
		if strings.Contains(resultSlice[0], "success") {
			return nil
		} else {
			log.Debug("Issue command: " + commandBuffer.String())
			log.Error("Fail to delete volume with error: " + resultSlice[0])
			return errors.New(resultSlice[0])
		}
	}
}