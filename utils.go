package main

import "github.com/mr-tron/base58"

// const (
//
//	WavesId          byte = 0x57 // blockchain ID
//	RawAddressPrefix byte = 0x01 // address first byte
//
// )

// var addressPrefix = []byte{RawAddressPrefix, WavesId}

func toBase58String(addr []byte) string {
	return base58.Encode(addr)
}

//func wavesHash(data []byte) []byte {
//	blakeHash, err := blake2b.New256(nil)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	blakeHash.Write(data)
//
//	blakeSum := blakeHash.Sum(nil)
//	keccakHash := sha3.NewLegacyKeccak256()
//
//	keccakHash.Write(blakeSum)
//
//	return keccakHash.Sum(nil)
//}

//func addressFromPubKey(data []byte) string {
//	// append generatorAddress hash (20 bytes) to the addressPrefix
//	generatorAddress := append(addressPrefix, wavesHash(data)[:20]...)
//
//	// append checksum (4 bytes) to the generatorAddress
//	checksum := wavesHash(generatorAddress)[:4]
//	generatorAddress = append(generatorAddress, checksum...)
//
//	// encode to base58
//	return base58.Encode(generatorAddress)
//}
