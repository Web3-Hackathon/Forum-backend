package main

import (
	"github.com/alxalx14/CryptoForum/Backend/Forum/database"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/server"
	"github.com/alxalx14/CryptoForum/Backend/Forum/server/sessions"
	"log"
)

func main() {
	var err error

	logger.Logger = logger.New("localhost", 5000, 1000)
	_, err = logger.Logger.Connect()
	if err != nil {
		log.Fatalf("[LOGGER] Could not connect to LogStash. Error: %s\n", err.Error())
	}

	go sessions.WatchDog()

	database.ConnectMongo("localhost:27017", "admin", "LoL187!!")
	database.ConnectMySQL("localhost:3306", "root", "LoL187!!")

	server.StartServer("0.0.0.0", 80)
}

//var client = rpc.New(rpc.DevNet_RPC)
//
//func getMetadata(mint solana.PublicKey) {
//	var programId = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
//
//	addr, _, err := solana.FindProgramAddress(
//		[][]byte{
//			[]byte("metadata"),
//			programId.Bytes(),
//			mint.Bytes(),
//		},
//		programId,
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	var info, _ = client.GetAccountInfo(context.TODO(), addr)
//	borshDec := bin.NewBorshDecoder(info.Value.Data.GetBinary())
//	var meta token_metadata.Metadata
//	err = borshDec.Decode(&meta)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(meta.Data.Uri)
//}

//func main() {
//	//var token = solana.MustPublicKeyFromBase58("H23Rn2MRapHxKD4aKfKtnLAwBvEbvEQgeAHQDuVHiEVj")
//	//getMetadata(token)
//
//	web3.GetAccountNFTs("8nBzhkZK1PQhqXXCZVWJyUXcTvKj9SVyJyU3uZnpVWD8")
//
//	//var sig, _ = hex.DecodeString("b92eaa734e8b1324fd5b41c0dc7e17d978d0fad0a734b31a10731145e453d122bb5035f5c17f86f68257a532c8d3476edbabf2b4d60cd59fbeb5e71a93c33103")
//	//
//	//var s = solana.SignatureFromBytes(sig)
//	//fmt.Println(s.Verify(pubKey, []byte("To avoid niggers, sign below to authenticate with white people")))
//
//}
