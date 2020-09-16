/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	networkextensionv1 "github.com/Tencent/bk-bcs/bcs-k8s/kubernetes/apis/networkextension/v1"
)

// SplitRegionLBID get region and lbid from regionLBID
func SplitRegionLBID(regionLBID string) (string, string, error) {
	strs := strings.Split(regionLBID, ":")
	if len(strs) == 1 {
		return "", strs[0], nil
	}
	if len(strs) == 2 {
		return strs[0], strs[1], nil
	}
	return "", "", fmt.Errorf("regionLBID %s invalid", regionLBID)
}

// GetListenerName generate listener name with lb id and port number
func GetListenerName(lbID string, port int) string {
	return lbID + "-" + strconv.Itoa(port)
}

// GetSegmentListenerName generate listener for port segment
func GetSegmentListenerName(lbID string, startPort, endPort int) string {
	return lbID + "-" + strconv.Itoa(startPort) + "-" + strconv.Itoa(endPort)
}

// GetPodIndex get pod index
func GetPodIndex(podName string) (int, error) {
	nameStrs := strings.Split(podName, "-")
	if len(nameStrs) < 2 {
		blog.Errorf("")
	}
	podNumberStr := nameStrs[len(nameStrs)-1]
	podIndex, err := strconv.Atoi(podNumberStr)
	if err != nil {
		blog.Errorf("get stateful set pod index failed from podName %s, err %s", podName, err.Error())
		return -1, fmt.Errorf("get stateful set pod index failed from podName %s, err %s", podName, err.Error())
	}
	return podIndex, nil
}

// GetDiffListeners get diff between two listener arrays
func GetDiffListeners(existedListeners, newListeners []networkextensionv1.Listener) (
	[]networkextensionv1.Listener, []networkextensionv1.Listener,
	[]networkextensionv1.Listener, []networkextensionv1.Listener) {

	existedListenerMap := make(map[string]networkextensionv1.Listener)
	for _, listener := range existedListeners {
		existedListenerMap[listener.GetName()] = listener
	}
	newListenerMap := make(map[string]networkextensionv1.Listener)
	for _, listener := range newListeners {
		newListenerMap[listener.GetName()] = listener
	}

	var adds []networkextensionv1.Listener
	var dels []networkextensionv1.Listener
	var olds []networkextensionv1.Listener
	var news []networkextensionv1.Listener

	for _, listener := range newListeners {
		existedListener, ok := existedListenerMap[listener.GetName()]
		if !ok {
			adds = append(adds, listener)
			continue
		}
		if !reflect.DeepEqual(listener.Spec, existedListener.Spec) {
			olds = append(olds, existedListener)
			news = append(news, listener)
			continue
		}
	}

	for _, listener := range existedListeners {
		_, ok := newListenerMap[listener.GetName()]
		if !ok {
			dels = append(dels, listener)
			continue
		}
	}
	return adds, dels, olds, news
}
