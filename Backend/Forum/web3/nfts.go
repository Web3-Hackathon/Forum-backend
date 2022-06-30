package web3

import (
	"context"
	"fmt"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// NFTInfo holds basic information about the NFT
// can be created using GetAccountNFTs
type NFTInfo struct {
	// Address is used to identify the token
	Address string
}

// Client is used to connect to the SOL network
var Client = rpc.New(rpc.MainNetBeta_RPC)

// TokenProgram is used to fetch NFTs of account
var TokenProgram = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

// MetadataProgram is used to fetch metadata of NFTs
var MetadataProgram = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

// GetAccountNFTs is used to fetch all NFTs owned by
// an account
func GetAccountNFTs(address string) []NFTInfo {
	var err error
	var pubKey solana.PublicKey
	var response *rpc.GetTokenAccountsResult

	response, err = Client.GetTokenAccountsByOwner(
		context.TODO(),
		pubKey,
		&rpc.GetTokenAccountsConfig{
			ProgramId: &TokenProgram,
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingJSONParsed,
		},
	)

	if err != nil {
		logger.Logf(logger.ERROR, "could not get NFTs of account. Error: %s", err.Error())
		return nil
	}

	//var nftInfo []NFTInfo

	for _, token := range response.Value {
		fmt.Println(token.Pubkey.String())
		//data, _ := token.Account.Data.MarshalJSON()
		//fmt.Println(string(data))
	}
	//nftInfo = append(nftInfo, NFTInfo{
	//	TokenAddress:  ,
	//	MintAuthority: "",
	//})

	//return nftInfo

	return nil
}
