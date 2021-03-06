/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package certs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/AISphere/ffdl-commons/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePublicPrivateKeyPair(t *testing.T) {

	id := time.Now().String()
	tmp, err := ioutil.TempDir("", id)
	assert.NoError(t, err)
	public := fmt.Sprintf("%s/public.pub", tmp)
	private := fmt.Sprintf("%s/private.pem", tmp)
	defer os.RemoveAll(tmp) //delete the folder once done
	assert.NoError(t, generatePublicPrivateKeyPair(public, private))
	if _, err := os.Stat(public); os.IsNotExist(err) {
		assert.Fail(t, "file should exist at the location %s", public)
	}
	if _, err := os.Stat(private); os.IsNotExist(err) {
		assert.Fail(t, "file should exist at the location %s", private)
	}

	//TODO load the private key and validate
}

func TestEqualityOfGeneratePublicPrivateKeyPair(t *testing.T) {

	assertForExistence := func(public, private string) {
		assert.NoError(t, generatePublicPrivateKeyPair(public, private))
		if _, err := os.Stat(public); os.IsNotExist(err) {
			assert.Fail(t, "file should exist at the location %s", public)
		}
		if _, err := os.Stat(private); os.IsNotExist(err) {
			assert.Fail(t, "file should exist at the location %s", private)
		}
	}

	assertForInequality := func(file1, file2 string) {
		fileContents1, err := ioutil.ReadFile(file1)
		assert.NoError(t, err)
		fileContents2, err := ioutil.ReadFile(file2)
		assert.NoError(t, err)
		assert.NotEqual(t, fileContents1, fileContents2)
	}

	id := time.Now().String()
	tmp, err := ioutil.TempDir("", id)
	assert.NoError(t, err)
	public1 := fmt.Sprintf("%s/public1.pub", tmp)
	private1 := fmt.Sprintf("%s/private1.pem", tmp)
	public2 := fmt.Sprintf("%s/public2.pub", tmp)
	private2 := fmt.Sprintf("%s/private2.pem", tmp)
	defer os.RemoveAll(tmp) //delete the folder once done
	assertForExistence(public1, private1)
	assertForExistence(public2, private2)
	assertForInequality(public1, public2)
	assertForInequality(private1, private2)

	//TODO load the private key and validate
}

func TestNeedsMountedSSHCerts(t *testing.T) {

	assert.True(t, NeedsMountedSSHCerts("tensorflow", "0.11-horovod"))
	assert.True(t, NeedsMountedSSHCerts("tensorflow", "1.3-py2-ddl"))
	assert.False(t, NeedsMountedSSHCerts("tensorflow", "1.4-py3"))
	assert.False(t, NeedsMountedSSHCerts("caffe2", "0.8"))
	assert.True(t, NeedsMountedSSHCerts("mxnet", "1.1.0"))
}

func TestFakeSecrets(t *testing.T) {

	name := "secretName"
	id := "trainingID"
	secret, err := GenerateSSHCertAsK8sSecret(name, id, "tensorflow", "0.11_horovod")
	assert.NoError(t, err)
	clientSet := fake.NewSimpleClientset(secret)
	recievedSecret, err := clientSet.Core().Secrets(config.GetLearnerNamespace()).Get(name, metav1.GetOptions{})
	assert.NoError(t, err)

	assert.EqualValues(t, recievedSecret.Labels["training_id"], id)
}
