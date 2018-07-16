/*
 * Copyright 2018 The ThunderDB Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package route

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/thunderdb/ThunderDB/conf"
	"gitlab.com/thunderdb/ThunderDB/proto"
	"gitlab.com/thunderdb/ThunderDB/utils/log"
)

const PubKeyStorePath = "./acl.keystore"

func TestIsPermitted(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	os.Remove(PubKeyStorePath)
	defer os.Remove(PubKeyStorePath)

	_, testFile, _, _ := runtime.Caller(0)
	confFile := filepath.Join(filepath.Dir(testFile), "../test/node_0/config.yaml")

	conf.GConf, _ = conf.LoadConfig(confFile)
	log.Debugf("GConf: %#v", conf.GConf)
	// reset the once
	Once = sync.Once{}
	InitKMS(PubKeyStorePath)

	Convey("test IsPermitted", t, func() {
		So(IsPermitted(conf.GConf.BP.NodeID, KayakCall), ShouldBeTrue)
		So(IsPermitted(proto.NodeID("0000"), KayakCall), ShouldBeFalse)
		So(IsPermitted(proto.NodeID("0000"), DHTFindNode), ShouldBeTrue)
		So(IsPermitted(proto.NodeID("0000"), RemoteFunc(9999)), ShouldBeFalse)
	})

	Convey("string RemoteFunc", t, func() {
		So(fmt.Sprintf("%s", DHTFindNode), ShouldContainSubstring, ".")
		So(fmt.Sprintf("%s", DHTFindNeighbor), ShouldContainSubstring, ".")
		So(fmt.Sprintf("%s", DHTPing), ShouldContainSubstring, ".")
		So(fmt.Sprintf("%s", MetricUploadMetrics), ShouldContainSubstring, ".")
		So(fmt.Sprintf("%s", KayakCall), ShouldContainSubstring, ".")
		So(fmt.Sprintf("%s", RemoteFunc(9999)), ShouldContainSubstring, "Unknown")
	})

}