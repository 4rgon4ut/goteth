package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func GetValIndexesFromJson(filePath string) ([]uint64, error) {

	var validatorIndex []uint64
	// open file and read all the indexes
	fbytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return validatorIndex, err
	}
	err = json.Unmarshal(fbytes, &validatorIndex)

	if err != nil {
		log.Errorf("Error unmarshalling val list: %s", err.Error())
	}

	log.Infof("Readed %d validators", len(validatorIndex))

	return validatorIndex, nil

}

func BoolToUint(input []bool) []uint64 {
	result := make([]uint64, len(input))

	for i, item := range input {
		if item {
			result[i] += 1
		}
	}
	return result
}

// in the case there is no pool
func DivideValidatorsBatches(input []uint64, workers int) []PoolKeys {

	result := make([]PoolKeys, 0)
	step := len(input) / workers

	includedIndex := 0
	for includedIndex < len(input) {
		endIndex := includedIndex + step
		if endIndex > len(input) { // to not overflow
			endIndex = len(input)
		}

		// from includedIndex to endIndex
		newBatch := PoolKeys{
			PoolName: "",
			ValIdxs:  input[includedIndex:endIndex],
			Pubkeys:  make([]string, 0),
		}
		result = append(result, newBatch)
		includedIndex = endIndex
	}
	return result
}

// From here we should obtain those validators that do not belong to any pool
func ObtainMissing(valLen int, poolVals [][]uint64) []uint64 {
	valList := make([]uint64, valLen) // initialized to 0, no need to track

	for _, poolArray := range poolVals {
		for _, item := range poolArray {
			valList[item] = 1 // it exists in the poolVals
		}
	}

	result := make([]uint64, 0)

	// track the validators that do not exist in the poolVals
	for i, item := range valList {
		if item == 0 {
			result = append(result, uint64(i))
		}
	}

	return result
}

func AddOthersPool(batches []PoolKeys, othervalList []uint64) []PoolKeys {

	for i, item := range batches {
		if item.PoolName == "others" {
			item.ValIdxs = append(item.ValIdxs, othervalList...)
			batches[i] = item
			return batches
		}
	}
	batches = append(batches, PoolKeys{
		PoolName: "others",
		ValIdxs:  othervalList,
		Pubkeys:  make([]string, 0),
	})
	return batches

}

func ReadCustomValidatorsFile(validatorKeysFile string) (validatorKeysByPool []PoolKeys, err error) {
	log.Info("Reading validator keys from: ", validatorKeysFile)
	validatorKeysByPool = make([]PoolKeys, 0)

	file, err := os.Open(validatorKeysFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip first line
		if line == "val_idx,pubkey,custom_pool" {
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			return validatorKeysByPool, errors.New("the format of the file is not the expected: f_val_idx, pubkey, pool_name")
		}

		// obtain three fields per line
		valIdx, err := strconv.Atoi(fields[0])
		pubkeyStr := strings.Trim(fields[1], "\"")
		poolName := fields[2]

		// check pubkey format
		pubkeyStr = strings.Replace(pubkeyStr, "\\x", "", -1)
		if !strings.HasPrefix(pubkeyStr, "0x") {
			pubkeyStr = "0x" + pubkeyStr
		}

		if len(pubkeyStr) != 98 {
			return validatorKeysByPool, errors.New(fmt.Sprintf("length of key for valIdx %d is incorrect: %d", valIdx, len(pubkeyStr)))
		}

		if err != nil {
			return validatorKeysByPool, errors.Wrap(err, fmt.Sprintf("could not parse valIdx: %d", valIdx))
		}

		found := false
		// look for which pool this line belongs to and append
		for i, item := range validatorKeysByPool {
			if poolName == item.PoolName {
				item.ValIdxs = append(item.ValIdxs, uint64(valIdx))
				item.Pubkeys = append(item.Pubkeys, pubkeyStr)
				validatorKeysByPool[i] = item
				found = true
				break
			}
		}
		if !found {
			valIdxs := make([]uint64, 0)
			valIdxs = append(valIdxs, uint64(valIdx))
			pubkeys := make([]string, 0)
			pubkeys = append(pubkeys, pubkeyStr)

			validatorKeysByPool = append(validatorKeysByPool, PoolKeys{
				PoolName: poolName,
				ValIdxs:  valIdxs,
				Pubkeys:  pubkeys,
			})

		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	log.Infof("Done reading from %s", validatorKeysFile)
	return validatorKeysByPool, nil
}

type PoolKeys struct {
	PoolName string
	ValIdxs  []uint64
	Pubkeys  []string
}
