// core_test.go: tests for core.go
//
// Copyright (c) 2013 CloudFlare, Inc.

package core

import (
	"bytes"
	"encoding/json"
	"os"
	"github.com/cloudflare/redoctober/keycache"
	"github.com/cloudflare/redoctober/passvault"
	"testing"
)

func TestCreate(t *testing.T) {
	createJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\"}")

	os.Remove("/tmp/db1.json")
	Init("/tmp/db1.json")

	respJson, err := Create(createJson)
	if err != nil {
		t.Fatalf("Error in creating account, ", err)
	}

	var s responseData
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in creating account, ", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in creating account, ", s.Status)
	}

	// check to see if creation can happen twice
	respJson, err = Create(createJson)
	if err != nil {
		t.Fatalf("Error in creating account when one exists, ", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in creating account when one exists, ", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in creating account when one exists, ", s.Status)
	}

	os.Remove("/tmp/db1.json")
}

func TestSummary(t *testing.T) {
	createJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\"}")
	delegateJson := []byte("{\"Name\":\"Bob\",\"Password\":\"Rob\",\"Time\":\"2h\",\"Uses\":1}")
	os.Remove("/tmp/db1.json")

	// check for summary of uninitialized vault
	respJson, err := Summary(createJson)
	if err != nil {
		t.Fatalf("Error in summary of account with no vault,", err)
	}
	var s summaryData
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in summary of account with no vault,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in summary of account with no vault, ", s.Status)
	}

	Init("/tmp/db1.json")

	// check for summary of initialized vault
	respJson, err = Create(createJson)
	if err != nil {
		t.Fatalf("Error in creating account, ", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in creating account, ", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in creating account, ", s.Status)
	}

	respJson, err = Summary(createJson)
	if err != nil {
		t.Fatalf("Error in summary of account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in summary of account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in summary of account, ", s.Status)
	}

	data, ok := s.All["Alice"]
	if !ok {
		t.Fatalf("Error in summary of account, record missing ")
	}
	if data.Admin != true {
		t.Fatalf("Error in summary of account, record incorrect ")
	}
	if data.Type != passvault.RSARecord {
		t.Fatalf("Error in summary of account, record incorrect ")
	}

	// check for summary of initialized vault with new member
	respJson, err = Delegate(delegateJson)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Summary(createJson)
	if err != nil {
		t.Fatalf("Error in summary of account with no vault,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in summary of account with no vault,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in summary of account with no vault, ", s.Status)
	}

	data, ok = s.All["Alice"]
	if !ok {
		t.Fatalf("Error in summary of account, record missing ")
	}
	if data.Admin != true {
		t.Fatalf("Error in summary of account, record missing ")
	}
	if data.Type != passvault.RSARecord {
		t.Fatalf("Error in summary of account, record missing ")
	}

	data, ok = s.All["Bob"]
	if !ok {
		t.Fatalf("Error in summary of account, record missing ")
	}
	if data.Admin != false {
		t.Fatalf("Error in summary of account, record missing ")
	}
	if data.Type != passvault.RSARecord {
		t.Fatalf("Error in summary of account, record missing ")
	}

	dataLive, ok := s.Live["Bob"]
	if !ok {
		t.Fatalf("Error in summary of account, record missing", keycache.UserKeys)
	}
	if dataLive.Admin != false {
		t.Fatalf("Error in summary of account, record missing ")
	}
	if dataLive.Type != passvault.RSARecord {
		t.Fatalf("Error in summary of account, record missing ")
	}

	keycache.FlushCache()

	os.Remove("/tmp/db1.json")
}

func TestPassword(t *testing.T) {
	createJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\"}")
	delegateJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"Time\":\"2h\",\"Uses\":1}")
	passwordJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"NewPassword\":\"Olleh\"}")
	delegateJson2 := []byte("{\"Name\":\"Alice\",\"Password\":\"Olleh\",\"Time\":\"2h\",\"Uses\":1}")
	passwordJson2 := []byte("{\"Name\":\"Alice\",\"Password\":\"Olleh\",\"NewPassword\":\"Hello\"}")
	os.Remove("/tmp/db1.json")

	Init("/tmp/db1.json")

	// check for summary of initialized vault with new member
	var s responseData
	respJson, err := Create(createJson)
	if err != nil {
		t.Fatalf("Error in creating account, ", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in creating account, ", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in creating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Password(passwordJson2)
	if err != nil {
		t.Fatalf("Error in password", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in password", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in password, ", s.Status)
	}

	respJson, err = Password(passwordJson)
	if err != nil {
		t.Fatalf("Error in password", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in password", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in password, ", s.Status)
	}

	respJson, err = Delegate(delegateJson)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson2)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Password(passwordJson2)
	if err != nil {
		t.Fatalf("Error in password", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in password", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in password, ", s.Status)
	}

	respJson, err = Delegate(delegateJson)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	keycache.FlushCache()

	os.Remove("/tmp/db1.json")
}

func TestEncryptDecrypt(t *testing.T) {
	summaryJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\"}")
	delegateJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"Time\":\"0s\",\"Uses\":0}")
	delegateJson2 := []byte("{\"Name\":\"Bob\",\"Password\":\"Hello\",\"Time\":\"0s\",\"Uses\":0}")
	delegateJson3 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"Time\":\"0s\",\"Uses\":0}")
	delegateJson4 := []byte("{\"Name\":\"Bob\",\"Password\":\"Hello\",\"Time\":\"10s\",\"Uses\":2}")
	delegateJson5 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"Time\":\"10s\",\"Uses\":2}")
	encryptJson := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"Minumum\":2,\"Owners\":[\"Alice\",\"Bob\",\"Carol\"],\"Data\":\"SGVsbG8gSmVsbG8=\"}")
	encryptJson2 := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"Minumum\":2,\"Owners\":[\"Alice\",\"Bob\",\"Carol\"],\"Data\":\"SGVsbG8gSmVsbG8=\"}")
	os.Remove("/tmp/db1.json")

	Init("/tmp/db1.json")

	// check for summary of initialized vault with new member
	var s responseData
	respJson, err := Create(delegateJson)
	if err != nil {
		t.Fatalf("Error in creating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in creating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in creating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson2)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson3)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	// check summary to see if none are delegated
	keycache.Refresh()
	respJson, err = Summary(summaryJson)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	var sum summaryData
	err = json.Unmarshal(respJson, &sum)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	if sum.Status != "ok" {
		t.Fatalf("Error in summary, ", sum.Status)
	}
	if len(sum.Live) != 0 {
		t.Fatalf("Error in summary, ", sum.Status)
	}

	// Encrypt with non-admin (fail)
	respJson, err = Encrypt(encryptJson)
	if err != nil {
		t.Fatalf("Error in encrypt,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in encrypt,", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in encrypt, ", s.Status)
	}

	// Encrypt
	respJson, err = Encrypt(encryptJson2)
	if err != nil {
		t.Fatalf("Error in encrypt,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in encrypt,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in encrypt, ", s.Status)
	}

	// decrypt file
	decryptJson, err := json.Marshal(decrypt{Name: "Alice", Password: "Hello", Data: s.Response})
	if err != nil {
		t.Fatalf("Error in marshalling decryption,", err)
	}

	respJson2, err := Decrypt(decryptJson)
	if err != nil {
		t.Fatalf("Error in decrypt,", err)
	}
	err = json.Unmarshal(respJson2, &s)
	if err != nil {
		t.Fatalf("Error in decrypt,", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in decrypt, ", s.Status)
	}

	// delegate two valid decryptors
	respJson, err = Delegate(delegateJson4)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson5)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	// verify the presence of the two delgations
	keycache.Refresh()
	var sum2 summaryData
	respJson, err = Summary(summaryJson)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	err = json.Unmarshal(respJson, &sum2)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	if sum2.Status != "ok" {
		t.Fatalf("Error in summary, ", sum2.Status)
	}
	if len(sum2.Live) != 2 {
		t.Fatalf("Error in summary, ", sum2.Live)
	}

	respJson2, err = Decrypt(decryptJson)
	if err != nil {
		t.Fatalf("Error in decrypt,", err)
	}
	err = json.Unmarshal(respJson2, &s)
	if err != nil {
		t.Fatalf("Error in decrypt,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in decrypt, ", s.Status)
	}
	if string(s.Response) != "Hello Jello" {
		t.Fatalf("Error in decrypt, ", string(s.Response))
	}

	keycache.FlushCache()

	os.Remove("/tmp/db1.json")
}

func TestModify(t *testing.T) {
	summaryJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\"}")
	summaryJson2 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\"}")
	delegateJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"Time\":\"0s\",\"Uses\":0}")
	delegateJson2 := []byte("{\"Name\":\"Bob\",\"Password\":\"Hello\",\"Time\":\"0s\",\"Uses\":0}")
	delegateJson3 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"Time\":\"0s\",\"Uses\":0}")
	modifyJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"ToModify\":\"Alice\",\"Command\":\"admin\"}")
	modifyJson2 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"ToModify\":\"Alice\",\"Command\":\"revoke\"}")
	modifyJson3 := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"ToModify\":\"Carol\",\"Command\":\"admin\"}")
	modifyJson4 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"ToModify\":\"Alice\",\"Command\":\"revoke\"}")
	modifyJson5 := []byte("{\"Name\":\"Carol\",\"Password\":\"Hello\",\"ToModify\":\"Alice\",\"Command\":\"delete\"}")

	os.Remove("/tmp/db1.json")
	Init("/tmp/db1.json")

	// check for summary of initialized vault with new member
	var s responseData
	respJson, err := Create(delegateJson)
	if err != nil {
		t.Fatalf("Error in creating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in creating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in creating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson2)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson3)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	// check summary to see if none are delegated
	keycache.Refresh()
	respJson, err = Summary(summaryJson)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	var sum summaryData
	err = json.Unmarshal(respJson, &sum)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	if sum.Status != "ok" {
		t.Fatalf("Error in summary, ", sum.Status)
	}
	if len(sum.Live) != 0 {
		t.Fatalf("Error in summary, ", sum.Status)
	}

	// Modify from non-admin (fail)
	respJson, err = Modify(modifyJson)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in modify, ", s.Status)
	}

	// Modify self from admin (fail)
	respJson, err = Modify(modifyJson2)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	if s.Status == "ok" {
		t.Fatalf("Error in modify, ", s.Status)
	}

	// Modify admin from admin
	respJson, err = Modify(modifyJson3)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in modify, ", s.Status)
	}

	respJson, err = Summary(summaryJson)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	err = json.Unmarshal(respJson, &sum)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	if sum.Status != "ok" {
		t.Fatalf("Error in summary, ", sum.Status)
	}
	if sum.All["Carol"].Admin != true {
		t.Fatalf("Error in summary, ", sum.All)
	}

	// Revoke admin from admin
	respJson, err = Modify(modifyJson4)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in modify, ", s.Status)
	}

	respJson, err = Summary(summaryJson2)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	err = json.Unmarshal(respJson, &sum)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	if sum.Status != "ok" {
		t.Fatalf("Error in summary, ", sum.Status)
	}
	if sum.All["Alice"].Admin == true {
		t.Fatalf("Error in summary, ", sum.All)
	}

	// Delete from admin
	respJson, err = Modify(modifyJson5)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in modify,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in modify, ", s.Status)
	}

	var sum3 summaryData
	respJson, err = Summary(summaryJson2)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	err = json.Unmarshal(respJson, &sum3)
	if err != nil {
		t.Fatalf("Error in summary,", err)
	}
	if sum3.Status != "ok" {
		t.Fatalf("Error in summary, ", sum3.Status)
	}
	if len(sum3.All) != 2 {
		t.Fatalf("Error in summary, ", sum3.All)
	}

	keycache.FlushCache()

	os.Remove("/tmp/db1.json")
}

func TestStatic(t *testing.T) {
	expected := []byte("this is a secret string, shhhhhhhh. don't tell anyone, ok? whatever you do, don't tell the cops\n")
	delegateJson2 := []byte("{\"Name\":\"Bob\",\"Password\":\"Rob\",\"Time\":\"1m\",\"Uses\":5}")
	delegateJson3 := []byte("{\"Name\":\"Carol\",\"Password\":\"~~~222\",\"Time\":\"1h\",\"Uses\":30}")
	decryptJson := []byte("{\"Name\":\"Alice\",\"Password\":\"Hello\",\"Data\":\"eyJWZXJzaW9uIjoxLCJWYXVsdElkIjo1Mjk4NTM4OTU3NzU0NTQwODEwLCJLZXlTZXQiOlt7Ik5hbWUiOlsiQm9iIiwiQWxpY2UiXSwiS2V5IjoiMmozdEkyUEJGQmJ3RlgwQmxRZFV1QT09In0seyJOYW1lIjpbIkJvYiIsIkNhcm9sIl0sIktleSI6InlMU1NCL1U2KzVyYzFFK2dqR1hUNHc9PSJ9LHsiTmFtZSI6WyJBbGljZSIsIkJvYiJdLCJLZXkiOiJERGxXSEY3c3p6SVN1WFdhRXo4bGxRPT0ifSx7Ik5hbWUiOlsiQWxpY2UiLCJDYXJvbCJdLCJLZXkiOiJUa0ExM2FQcllGTmJ2ZUlibDBxZHd3PT0ifSx7Ik5hbWUiOlsiQ2Fyb2wiLCJCb2IiXSwiS2V5IjoiS1htMHVPYm1SSjJadlNZRVdQd2syQT09In0seyJOYW1lIjpbIkNhcm9sIiwiQWxpY2UiXSwiS2V5IjoiTDljK1BxdHhQaDl5NmFwUnZ0YkNRdz09In1dLCJLZXlTZXRSU0EiOnsiQWxpY2UiOnsiS2V5IjoiZmo1bXFucTd5NUtDYWZDR1QxSTUxeEk1SnNYNzQ2KzlUVFNzcC84eWJmM2laamhGelNsd1AzYU5tc094M1NVS1RtWmxmcytiK01lRDRlS0oydUtCRnpRSEFJUE8wZndvaUNES0hoS0g2S3NvbE5xNCtqZ3BrTEFNT0xzUUdzOGc2QmhKeTZiQ0ZSalpWYzNJZGxRQUJQTTZQa1RidVN2S2huOWF0REZ3WlFENVRKQklTaTdkMmh3NExyYWR0TElUYk5xd2lGTVRRUTkrcHNYenlhdlk4SDNMTkhLR2dmNU9kN0lwdGhFUVBDSGk0bnc3WC9ZVlJURU1mb0lWY01jS093WWpsQzQ1L1ZKRUhLOVp5M0RTaUxCemptcjU3WU5JVmp3OFlaWTVER0JXcWJndTUxUlViSWNycXlMcGhCaFhvQlJ1NFIreXJoeWdCTldidmtraWZBPT0ifSwiQm9iIjp7IktleSI6IkhUWWlaMThzZjcyMWNBTjFMUk5rSi8rTDRBS1dpbE1ya015TmlCaldjbDlIUlRWUE5YSVRxUUJYZDBmQmdnR05QaVpyNlZRVHlTSzRaRnZKS0dER2l6MTdUZS9Ub0RuOFlrL0I5Y3FNc041ZkhvUXRYdmw4SVpvMndpb0E2N2NjQUoxZ0hNTU5wUHlMZEY0M1NRcWdJK1hhUTJsTVNZTE1meHhEbUJCT1ExU1dBdG8wQkRSZG5zcXB3VXdJUEtROVkzLzFvc21yakxtSm9BQzNNUHBsZXhZV2hleE53SnRTZCttRmRWWjNRZTR4OVJzUkhjTi9teWloT3QvNjdWNjBxenMxM0YwUlprTVNEemo1RGRnKzFLVk5KWlk5ZG1vbFBOa0Faajh6MjBMOXV6cGF0clRZVFI2QThxL3NSbitpbk83WlFWUTAwWE82cTZsWVlRenhudz09In0sIkNhcm9sIjp7IktleSI6Ikl0cnZTMDJuU2ZiY0EyZmwxTDFpNjF4cVBFREtSZHNyWWUzK1VDYmtUK2lwaGVpUVJQU3Vpa2J6ZVYya3NobjR5SkRla3U1Ym1UTnFXOEhTR3RVN0dUZ0NvSVdWOFdtRWY0dzZvdnpTaFBidStWckladlJ6M3dqaDJvWUhUL2d0UFZBUW5CYS83MUZlb0JOeHk1bC9oQmNVbUJreTQzajgzTWx0Mis4UVp4NlBFVURtcGFQUWVtVmg5OStDMjBuUXRrQVVGZU1jMkdlNHk3UmxIU3h0ZkFCdndsWHgxTnpDRDQwbnlKZkYxU2pWL2ZaaC9FMkFsNFRhdng2RE9Ka1lHb0oybXA3WEJ2WDBJRjJ0cDhUM1U1VnBuVGVrL1d1TnJMTDl6Ny9qcXpXaDg3bFo1S2hlV1hoR2tVMUJOSDRsZklqNDNwRGtTeTUwYUR2UzB6WWZIUT09In19LCJJViI6IjU4cjlNejhlMDZtSXRCRzluU1YvMFE9PSIsIkRhdGEiOiJRRTlaaGNHWE5YYXVVZE1rMDRiaVVHeTFTb1A1SDJuRi9qMkpqaWlWRktQZElkUnAvR2MrQVp2VUk5bjIyWk00cSt6RGlKejdxdks0YkthUHBYaFRtR1AwWGhlYUZVdWtlVk5TOVNUTW9UYk5jWS9adFZPejZoaXpVUEY3Z1NxMzg4UVBVc1QrQXhtbDNyRVVUV09obnc9PSIsIlNpZ25hdHVyZSI6IlZWaVJvYUN4cUtBdk54eU14RGhydFlvMENaQT0ifQ==\"}")
	diskVault := []byte("{\"Version\":1,\"VaultId\":5298538957754540810,\"HmacKey\":\"Qugc5ZQ0vC7KQSgmDHTVgQ==\",\"Passwords\":{\"Alice\":{\"Type\":\"RSA\",\"PasswordSalt\":\"HrVbQ4JvEdORxWN6FrSy4w==\",\"HashedPassword\":\"nanadB0t/EmVuHyRUNtfyQ==\",\"KeySalt\":\"u0cwMmHikadpxm9wClrf3g==\",\"AESKey\":\"pOcxCYMk0l+kaEM5IHolfA==\",\"RSAKey\":{\"RSAExp\":\"yDBotK0XkaDk6rt35ciBnlyPGSXjp9ypdTH2j89CybXe2ReF6xLVZ10CCoz91UUFpbiQi2tVWFS8V7lxUehx7C6HI30Zr0iwJMXk8sgRKs2Ee3rAGCrM9vvQMO8ApKwe4kvB+PrFgnhIEgZI7zyPJPcdZnxbcVFyUC3uIivFS4Bv3jVBnrkAx71keQqrkKzu3jTquCIS3rTEm2hrgsZ3n1t/4BADvUpcMpII2Phv2Senonp1EMz9sRFsR4w9HVep8AgfUGuPTcQ/uv9R16xiHvlN80rOtWzVLpEruzGw/JTvlsshFBU5SY/zthGl9TwxWz1yZpVHYhIWhbxw34CXYczB6q19bdEPkJJldi/1coI=\",\"RSAExpIV\":\"RI9B8nzwW/CXiIU4lSYRWg==\",\"RSAPrimeP\":\"yB3TMdRhWLiZ/ayPlu/iLDHxWsuMi9pA8Ctps7WZFxVsxYEzy/s0Otzsbtgay8of72xVkO64MaudXXRjj26kVSQhS8WhPmv7xDO5ba3SsTffCA99aSr3MH+JmUoL+EDjYviUf5F1DSZniv+Ae+6x9AMbbRQqRvmdH/INW76rFbLX4VtsMpgVkAhADwCUSfNS\",\"RSAPrimePIV\":\"ZE8xe7zVqz4fTN1XWvD1Tg==\",\"RSAPrimeQ\":\"pq7rEAmXoiMWhuEpTy2pQiheLHuzcGeGm1IsgtP7TRIpHaumBAjasxhMY95ODspzehfHp8RPb3pz550g+EXpRP5bfmZkiPMGWKsFUWY7h51hm2Yg6t7yOSXxQvzseWDUjJqBXIG56su0ItD/7NT9YPnhOWAcLaV+L/D3dLdUOP3sgrfp2fcz5IWgcAHUKwLj\",\"RSAPrimeQIV\":\"po4GB9LNDJEFwVo7pd7M+w==\",\"RSAPublic\":{\"N\":23149224936745228096494487629641309516714939646787578848987096719230079441991920927625328021835418605829165905957281054739040782220661415211495551591590967025905842090248342119030961900558364831442589592519977574953153506098793781379522492850670706181134690851176026007414311463584686376540335844073069224992659463459793269384975134553539299820716318523937519020382533972574258510063896159690996712751336541794097684994088922002126390397554509324964165816682788604802958437629743951677011277308533675371215221012180522157452368867794857042597831942805118364469257052754518983480732149823871291876054393236400292711237,\"E\":65537}},\"Admin\":true},\"Bob\":{\"Type\":\"RSA\",\"PasswordSalt\":\"yqWZFNuuCRMw+snL83FcWQ==\",\"HashedPassword\":\"jYLSUIdvw8UVXhxVSS5uKQ==\",\"KeySalt\":\"ZJDOGSv9yUlJ5+83+FV6Nw==\",\"AESKey\":\"9bj4qQDJKglD9eIm7MNeag==\",\"RSAKey\":{\"RSAExp\":\"rnmHX1FgAN0KFm89Uj+O4L7/njDlQPwGHSYjGKB7qMyNPGF4jKn9ez1LI9e9jI8uE475KBhoXRmIuHdap6HtD+sH+nHugDzD3np6Dbq1MM+19PW41n9xblwx6tsNwvufCYgb2qtZd0bLXemeKBszJg2UXlTeSHkXGmjY/VAbZYbUFUNKIykiJWnQ+3HlAo7UHjjKSDI6HtiMnODwYucKH9uYdmroM3DqRUP+j+AhMctyeyOt68q6RuVyubzG0PM1/T9QqPgTIzFvg9dxk+LmPYmlv4b7Euea7KQHww4kTUlpNYondRisuA436G7EfWIJFeRShFuSVk0GvmrgN5vK+/FkdqMuikPaPV1dITFzaCU=\",\"RSAExpIV\":\"eoIov2XYwK19tpFkZyTL+Q==\",\"RSAPrimeP\":\"x9VqeHfHZBMksWNz3Bq90KyLYO0+y9tbfH17uAxQHIZbhn6RUHMXJFs4/H+TnL9s7D05HdxKNCD3ilISkhLZ/DQVD+VMvSFQ/2DL/OoKNg4UyxeGpKefJxCU9LOeXGTfN+UpGHN4Myal3PL2yi0yWCVZvX+zPks05Nqmzkmtx1Yz0IqeUaR+I1jw2y0xy+nc\",\"RSAPrimePIV\":\"h+J0h35kSaNFjBFyejmqSQ==\",\"RSAPrimeQ\":\"CUGaukdC8slcy1Qm6kBw8QYhVenJWFylPqJTvf03ATkkaxQJ8ZSK/RodXFiKCztRSx/HVw0LoGGx9RiavSxl1I+NPIiGEBXiYnLCwWjIbohEIviX0XrKEegECKZtEPxlkDGI8C5ScmUiUOoCUFODqPi1ymEPPIpqj6NZu0Q3lOXseZ7vHbCwiO7Nxznoed2V\",\"RSAPrimeQIV\":\"fglAccK/JKbbP6FKd3MoPA==\",\"RSAPublic\":{\"N\":28335031342838743289078885439639803673108967200362163631216970159249785951860249148476503896024171065194110189927394386737974335124568755441420594193064949823336810514078294805084876384362160530316962043860952904024334641317905429962254720146672817625880231384279651687169091903559644110839987019710966818515746956556387614880164417442090115347521394174077031531398826388780347541493782387177778507223791525305993050791342196579264299166530031123364885677134064847664024722422373550942639559346558112737229502294633354781474882714474174085566901518771722057743396746938093859383451665388899451851738694979184054459471,\"E\":65537}},\"Admin\":false},\"Carol\":{\"Type\":\"RSA\",\"PasswordSalt\":\"kPCA68PNIi3qODVfBHonMA==\",\"HashedPassword\":\"L05ueGkNhqRI7xo7OaG5fQ==\",\"KeySalt\":\"bUXHAjmdtpcYMNN/YzHvhA==\",\"AESKey\":\"34sh2HYt11JMd63CTC2g3Q==\",\"RSAKey\":{\"RSAExp\":\"blJ6FvmAC6ZegcR6ITdg8WDSMljcikrUp3dRGbRDgKIfK0sx7b9i/kDtp4uD7/3IuTpD/qq09k07PO10T5R9NI2OdaEUGsoKJ4wqyzP6XzC9l/KeQmU7cMTh5OHO5UcqXWBn+g1INWaWI6DNAPFKK+4jcTyy6gTQQZQKSYRvRmgN6GkjOhbsdH0eM29eZScShJy2cGemZwbx8g2x7+cebALJmnJxycq29uBZaNBq4gV2/oNXxUhAogJY8SVhgPzCgig2MAJBLK/PzzxJ2LgnBLosZkUXo4vLPKSELX0qXSkSOoC/qPXvcQGawsknqUkSaaDwB1T1MYHmcJS/wpiIRM89Qru0Oy8sy8NRApk3ef0=\",\"RSAExpIV\":\"pQnJ6tvdJkPmMkqFjwJg5g==\",\"RSAPrimeP\":\"xTwgxhDvEvMCn8W6Mqv4ogRnBXFTLbuLP8YC2exQ+RwKmUbNsjB1eBlHyYRtKwc7qoDLf/zqpMZk3yPPcQDP+szAS5mmhCgpd+ePee8vvwVR5DImDuCquMGrPNEgpg3LL1aBHMF66pfnYGob+P6GZYzJP3mhlewUlh86NDzP8YkoVlrNRZrarB1/ZmA+w2Y2\",\"RSAPrimePIV\":\"BYynxyFw7WxrpSKJ3ogwOw==\",\"RSAPrimeQ\":\"wrSXyN/gpDTvK/BMeH+04uhzGUFVbBrfo47L3i+E1QsQndJDy2yyyE1D26mFcvyiefObaD097M3ruR9Wz7WfjxmyawG8N2W/BgtHZz/Ds2ThN/T9t1XqskxnIV8l2eq9LL7SqjAyp8jUGMNVS4WODDtDngLIR6OeehfyHpgwVZRuqQMxIT3wA78SwDcWo41s\",\"RSAPrimeQIV\":\"nh42rsJHRpRyrT7DUHCSGQ==\",\"RSAPublic\":{\"N\":21157349824302424474586534982259249211803834319040688464355493234801787797103302507532796077498450217821871970018694577875997137564403612290101696297388762661598985028358214691463545987525906601624475015641369758738594698734277968689933610815708041386873182005594893971935327515126641632178583839753604241275307504709356922001828087700047657269629468786192584360390570302228520684623917906322926285597325969441240603594319235541820769758760147548185811818605055779607442920451051925838835969679126714958513824527149123555098280031335447316080110515987423541389176918537439220114191601986031091845506241054078769049741,\"E\":65537}},\"Admin\":false}}}")

	file, err := os.Create("/tmp/db1.json")
	if err != nil {
		t.Fatalf("Error opening file,", err)
	}

	_, err = file.Write(diskVault)
	if err != nil {
		t.Fatalf("Error writing file,", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("Error closing file,", err)
	}

	Init("/tmp/db1.json")

	// check for summary of initialized vault with new member
	var s responseData
	respJson, err := Delegate(delegateJson2)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	respJson, err = Delegate(delegateJson3)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	err = json.Unmarshal(respJson, &s)
	if err != nil {
		t.Fatalf("Error in delegating account,", err)
	}
	if s.Status != "ok" {
		t.Fatalf("Error in delegating account, ", s.Status)
	}

	var r responseData
	respJson, err = Decrypt(decryptJson)
	if err != nil {
		t.Fatalf("Error in decrypt,", err)
	}
	err = json.Unmarshal(respJson, &r)
	if err != nil {
		t.Fatalf("Error in decrypt,", err)
	}
	if r.Status != "ok" {
		t.Fatalf("Error in summary, ", r.Status)
	}

	if bytes.Compare(expected, r.Response) != 0 {
		t.Fatalf("Error in summary, ", expected, r.Response)
	}

	keycache.FlushCache()

	os.Remove("/tmp/db1.json")
}
