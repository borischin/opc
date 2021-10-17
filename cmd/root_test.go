package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

const TEST_DATA_DIR = "../test/data"

func Test_InvalidArgument(t *testing.T) {
	test(t, []string{"-k"}, func(out []byte) {
		expectedMsg := "unknown shorthand flag: 'k' in -k"
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf(showDiff([]byte(expectedMsg), out))
		}
	})
}

func Test_ArgumentModuleRequired(t *testing.T) {
	test(t, []string{}, func(out []byte) {
		expectedMsg := "required flag(s) \"module\" not set"
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf(showDiff([]byte(expectedMsg), out))
		}
	})
}

func Test_ArgumentModuleFileNotFound(t *testing.T) {
	test(t, []string{"-m", "mmm"}, func(out []byte) {
		expectedMsg := "open mmm: no such file or directory"
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf(showDiff([]byte(expectedMsg), out))
		}
	})
}

func Test_ArgumentInvalidInput(t *testing.T) {
	test(t, []string{"-m", "mmm", "-i", "aaa"}, func(out []byte) {
		expectedMsg := "Invalid input(aaa)."
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf(showDiff([]byte(expectedMsg), out))
		}
	})
}

func Test_ArgumentInvalidOutputFormat(t *testing.T) {
	test(t, []string{"-m", "mmm", "-f", "aaa"}, func(out []byte) {
		expectedMsg := "Invalid outputFormat(aaa)."
		if !strings.Contains(string(out), expectedMsg) {
			t.Fatalf(showDiff([]byte(expectedMsg), out))
		}
	})
}

func Test_NoPolicyConfig(t *testing.T) {
	test(t, []string{"-m", getPath("no_policy_config.rego")}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("no_policy_config_exp.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_NoPolicyConfig_EnvFile(t *testing.T) {
	test(t, []string{"-m", getPath("no_policy_config.rego"), "-f", "env-file"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("no_policy_config_exp.env"))

		assertEqual(t, out, expectedMsg)
	})
}

func Test_SimplePolicyConfig_NoInput(t *testing.T) {
	test(t, []string{"-m", getPath("simple_policy_config.rego")}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("simple_policy_config_exp_no_input.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_SimplePolicyConfig_InputQa(t *testing.T) {
	test(t, []string{"-m", getPath("simple_policy_config.rego"), "-i", "env=QA"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("simple_policy_config_exp_qa.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_SimplePolicyConfig_InputProd(t *testing.T) {
	test(t, []string{"-m", getPath("simple_policy_config.rego"), "-i", "env=PROD"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("simple_policy_config_exp_prod.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_ComplicatedPolicyConfig_NoInput(t *testing.T) {
	test(t, []string{"-m", getPath("complicated_policy_config.rego")}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("complicated_policy_config_exp_no_input.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_ComplicatedPolicyConfig_InputQA_TW(t *testing.T) {
	test(t, []string{"-m", getPath("complicated_policy_config.rego"), "-i", "env=QA", "-i", "market=TW"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("complicated_policy_config_exp_qa_tw.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_ComplicatedPolicyConfig_InputQA_US(t *testing.T) {
	test(t, []string{"-m", getPath("complicated_policy_config.rego"), "-i", "env=QA", "-i", "market=US"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("complicated_policy_config_exp_qa_us.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_ComplicatedPolicyConfig_InputPROD_TW(t *testing.T) {
	test(t, []string{"-m", getPath("complicated_policy_config.rego"), "-i", "env=PROD", "-i", "market=TW"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("complicated_policy_config_exp_prod_tw.json"))

		assertJsonEqual(t, out, expectedMsg)
	})
}

func Test_ComplicatedPolicyConfig_InputPROD_US(t *testing.T) {
	test(t, []string{"-m", getPath("complicated_policy_config.rego"), "-i", "env=PROD", "-i", "market=US"}, func(out []byte) {
		expectedMsg, _ := ioutil.ReadFile(getPath("complicated_policy_config_exp_prod_us.json"))

		assertJsonEqual(t, out, expectedMsg)
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

func assertJsonEqual(t *testing.T, a, b []byte) {
	j1, err := jsonPretty(a, "    ")
	if err != nil {
		t.Fatalf("\nError: jsonPretty failed, %s\n%s", err, a)
		return
	}

	j2, err := jsonPretty(b, "    ")
	if err != nil {
		t.Fatalf("\nError: jsonPretty failed, %s\n%s", err, b)
		return
	}

	j1 = strings.TrimSpace(j1)
	j2 = strings.TrimSpace(j2)

	if j1 != j2 {
		t.Fatalf(showDiff([]byte(j1), []byte(j2)))
	}
	return
}

func assertEqual(t *testing.T, a, b []byte) {
	s1 := strings.TrimSpace(string(a))
	s2 := strings.TrimSpace(string(b))

	if s1 != s2 {
		t.Fatalf(showDiff([]byte(s1), []byte(s2)))
	}
	return
}

func getPath(file string) string { return path.Join(TEST_DATA_DIR, file) }

func showDiff(a, b []byte) string {
	return fmt.Sprintf("\n[Expected msg contains]\n%s\n-----------------\n[Actual msg]\n%s\n-----------------", a, b)
}
