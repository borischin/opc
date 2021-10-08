package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func Test_ModuleRequired(t *testing.T) {
	test(t, []string{}, func(out []byte) {
		expectedMsg := "required flag(s) \"module\" not set"
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf("expected error msg contains: %s", expectedMsg)
		}
	})
}

func Test_ModuleFileNotFound(t *testing.T) {
	test(t, []string{"-m", "aaa"}, func(out []byte) {
		expectedMsg := "open aaa: no such file or directory"
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf("\n[Expected msg contains]\n%s\n-----------------\n[Actual msg]\n%s\n-----------------", expectedMsg, out)
		}
	})
}

func Test_NoPolicyConfig(t *testing.T) {
	test(t, []string{"-m", "../test/data/no_policy_config.rego"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile("../test/data/no_policy_config_exp.json")
		isEqual, _ := jsonEqual(out, expectedMsg)

		if isEqual == false {
			prettyJson, _ := jsonPretty(out, "    ")
			t.Fatalf("\n[Expected msg contains]\n%s\n-----------------\n[Actual msg]\n%s\n-----------------", expectedMsg, prettyJson)
		}
	})
}

func test(t *testing.T, args []string, expect func(out []byte)) {
	b := bytes.NewBufferString("")
	cmd := GetRootCmd()
	cmd.SetErr(b)
	cmd.SetOut(b)
	cmd.SetOutput(b)
	cmd.SetArgs(args)
	cmd.Execute()

	out, err := ioutil.ReadAll(b)

	if err != nil {
		t.Fatal(err)
	}

	expect(out)
}

func jsonPretty(input []byte, indent string) (string, error) {
	dst := &bytes.Buffer{}
	err := json.Indent(dst, []byte(input), "", indent)
	if err != nil {
		return "", err
	}
	return dst.String(), nil
}

func jsonEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
