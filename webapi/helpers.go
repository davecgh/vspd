// Copyright (c) 2021-2023 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package webapi

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/decred/dcrd/blockchain/stake/v5"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrutil/v4"
	dcrdtypes "github.com/decred/dcrd/rpc/jsonrpc/types/v4"
	"github.com/decred/dcrd/wire"
	"github.com/decred/vspd/database"
)

func currentVoteVersion(params *chaincfg.Params) uint32 {
	var latestVersion uint32
	for version := range params.Deployments {
		if latestVersion < version {
			latestVersion = version
		}
	}
	return latestVersion
}

// validConsensusVoteChoices returns an error if provided vote choices are not
// valid for the most recent consensus agendas.
func validConsensusVoteChoices(params *chaincfg.Params, voteVersion uint32, voteChoices map[string]string) error {

agendaLoop:
	for agenda, choice := range voteChoices {
		// Does the agenda exist?
		for _, v := range params.Deployments[voteVersion] {
			if v.Vote.Id == agenda {
				// Agenda exists - does the vote choice exist?
				for _, c := range v.Vote.Choices {
					if c.Id == choice {
						// Valid agenda and choice combo! Check the next one...
						continue agendaLoop
					}
				}
				return fmt.Errorf("choice %q not found for agenda %q", choice, agenda)
			}

		}
		return fmt.Errorf("agenda %q not found for vote version %d", agenda, voteVersion)
	}

	return nil
}

func validTreasuryPolicy(policy map[string]string) error {
	for key, choice := range policy {
		pikey, err := hex.DecodeString(key)
		if err != nil {
			return fmt.Errorf("error decoding treasury key %q: %w", key, err)
		}
		if len(pikey) != secp256k1.PubKeyBytesLenCompressed {
			return fmt.Errorf("treasury key %q is not 33 bytes", key)
		}

		err = validPolicyOption(choice)
		if err != nil {
			return err
		}
	}

	return nil
}

func validTSpendPolicy(policy map[string]string) error {
	for hash, choice := range policy {
		if len(hash) != chainhash.MaxHashStringSize {
			return fmt.Errorf("wrong tspend hash length, expected %d got %d",
				chainhash.MaxHashStringSize, len(hash))
		}

		_, err := chainhash.NewHashFromStr(hash)
		if err != nil {
			return fmt.Errorf("error decoding tspend hash %q: %w", hash, err)
		}

		err = validPolicyOption(choice)
		if err != nil {
			return err
		}
	}

	return nil
}

// validPolicyOption checks that policy is one of the valid values accepted by
// dcrwallet RPCs. Invalid values return an error.
func validPolicyOption(policy string) error {
	switch policy {
	case "yes", "no", "abstain", "invalid", "":
		return nil
	default:
		return fmt.Errorf("%q is not a valid policy option", policy)
	}
}

func validateSignature(hash, commitmentAddress, signature, message string,
	db *database.VspDatabase, params *chaincfg.Params) error {

	firstErr := dcrutil.VerifyMessage(commitmentAddress, signature, message, params)
	if firstErr != nil {
		// Don't return an error straight away if sig validation fails -
		// first check if we have an alternate sign address for this ticket.
		altSigData, err := db.AltSignAddrData(hash)
		if err != nil {
			return fmt.Errorf("db.AltSignAddrData failed: %w", err)
		}

		// If we have no alternate sign address, or if validating with the
		// alt sign addr fails, return an error to the client.
		if altSigData == nil ||
			dcrutil.VerifyMessage(altSigData.AltSignAddr, signature, message, params) != nil {
			return fmt.Errorf("bad signature")
		}

	}
	return nil
}

func decodeTransaction(txHex string) (*wire.MsgTx, error) {
	msgHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	msgTx := wire.NewMsgTx()
	if err = msgTx.FromBytes(msgHex); err != nil {
		return nil, err
	}
	return msgTx, nil
}

func isValidTicket(tx *wire.MsgTx) error {
	if !stake.IsSStx(tx) {
		return errors.New("invalid transaction - not sstx")
	}
	if len(tx.TxOut) != 3 {
		return fmt.Errorf("invalid transaction - expected 3 outputs, got %d", len(tx.TxOut))
	}

	return nil
}

// validateTicketHash ensures the provided ticket hash is a valid ticket hash.
// A ticket hash should be 64 chars (MaxHashStringSize) and should parse into
// a chainhash.Hash without error.
func validateTicketHash(hash string) error {
	if len(hash) != chainhash.MaxHashStringSize {
		return fmt.Errorf("incorrect hash length: got %d, expected %d", len(hash), chainhash.MaxHashStringSize)

	}
	_, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return fmt.Errorf("invalid hash: %w", err)

	}

	return nil
}

// canTicketVote checks determines whether a ticket is able to vote at some
// point in the future by checking that it is currently either immature or live.
func canTicketVote(rawTx *dcrdtypes.TxRawResult, dcrdClient Node, netParams *chaincfg.Params) (bool, error) {

	// Tickets which have more than (TicketMaturity+TicketExpiry+1)
	// confirmations are too old to vote.
	if rawTx.Confirmations > int64(uint32(netParams.TicketMaturity)+netParams.TicketExpiry)+1 {
		return false, nil
	}

	// If ticket is currently immature, it will be able to vote in future.
	if rawTx.Confirmations <= int64(netParams.TicketMaturity) {
		return true, nil
	}

	// If ticket is currently live, it will be able to vote in future.
	live, err := dcrdClient.ExistsLiveTicket(rawTx.Txid)
	if err != nil {
		return false, fmt.Errorf("dcrd.ExistsLiveTicket error: %w", err)
	}

	return live, nil
}
