package fasdk

import (
	"bytes"
	"fmt"
	"encoding/pem"
	"encoding/base64"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/hyperledger/fabric/common/configtx"
	"github.com/hyperledger/fabric/protos/common"
	cm "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"crypto/x509"
	"github.com/hyperledger/fabric/core/ledger/util"
	"github.com/jonluo94/baasmanager/baas-core/common/json"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
)

type FabricBlock struct {
	Number uint64 `json:"number"`
	CurrentBlockHash string `json:"currentBlockHash"`
	PreviousBlockHash string `json:"previousBlockHash"`
	Transactions []FabricTransaction `json:"transactions"`
}

type FabricTransaction struct {
	No int64 `json:"no"`
	Status string `json:"status"`
	Txid string `json:"txid"`
	Channel string `json:"channel"`
	Type string `json:"type"`
	Subject string `json:"subject"`
	Config string `json:"config"`
	Actions []FabricTransactionAction `json:"actions"`
}

type FabricTransactionAction struct {
	Endorsers []string `json:"endorsers"`
	RWSet  []FabricTransactionActionRWSet `json:"rwSet"`
}

type FabricTransactionActionRWSet struct {
	Cc string `json:"cc"`
	RSet []string `json:"rSet"`
	WSet []string `json:"wSet"`
}


func blockParse(block *cm.Block ) *common.Block{
	cmBlock := new(common.Block)
	cmBlock.Data = &common.BlockData{
		Data:block.Data.Data,
	}
	cmBlock.Header= &common.BlockHeader{
		Number: block.Header.Number,
		PreviousHash: block.Header.PreviousHash,
		DataHash: block.Header.DataHash,
	}
	cmBlock.Metadata =&common.BlockMetadata{
		Metadata: block.Metadata.Metadata,
	}
	return cmBlock
}

func decodeSerializedIdentity(creator []byte) (string, error) {
	certStart := bytes.IndexAny(creator, "-----BEGIN")
	if certStart == -1 {
		return "", fmt.Errorf("No certificate found")
	}
	certText := creator[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return "", fmt.Errorf("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return "", err
	}
	uname := cert.Subject.CommonName
	return uname, nil
}
func parseBlock(block *common.Block) (*FabricBlock,error) {
	var err error
	faBlock := new(FabricBlock)
	trans := make([]FabricTransaction,0)
	// Handle header
	faBlock.Number = block.GetHeader().Number
	faBlock.CurrentBlockHash = base64.StdEncoding.EncodeToString(block.GetHeader().DataHash)
	faBlock.PreviousBlockHash = base64.StdEncoding.EncodeToString(block.GetHeader().PreviousHash)
	// Handle transaction
	var tranNo int64 = -1
	txsFilter := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
	if len(txsFilter) == 0 {
		txsFilter = util.NewTxValidationFlags(len(block.Data.Data))
		block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
	}

	for _, envBytes := range block.Data.Data {
		tran := FabricTransaction{}
		tranNo++
		tran.No = tranNo
		if txsFilter.IsInvalid(int(tranNo)) {
			tran.Status = "INVALID"
			continue
		} else {
			tran.Status = "VALID"
		}

		var env *common.Envelope
		if env, err = utils.GetEnvelopeFromBlock(envBytes); err != nil {
			return nil,err
		}

		var payload *common.Payload
		if payload, err = utils.GetPayload(env); err != nil {
			return nil,err
		}

		var chdr *common.ChannelHeader
		chdr, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if  err != nil {
			return nil,err
		}
		tran.Txid = chdr.TxId
		tran.Channel = chdr.ChannelId

		var shdr *common.SignatureHeader
		shdr, err = utils.GetSignatureHeader(payload.Header.SignatureHeader)
		if err != nil {
			return nil,err
		}

		var subject string
		subject, err = decodeSerializedIdentity(shdr.Creator)
		if  err != nil {
			return nil,err
		}
		tran.Subject = subject
		if common.HeaderType(chdr.Type) == common.HeaderType_CONFIG {
			tran.Type = "CONFIG"
			conf, err := parseConfig(payload)
			if err != nil {
				return nil,err
			}
			tran.Config = string(conf)
		} else if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {
			tran.Type = "ENDORSER_TRANSACTION"
			actions, err := parseEndorserTransaction(payload)
			if err != nil {
				return nil,err
			}
			tran.Actions = actions
		} else {
			tran.Type = "UNKNOWN"
		}

		trans = append(trans,tran)
	}
	faBlock.Transactions = trans
	return faBlock,nil
}

func parseEndorserTransaction(payload *common.Payload) ([]FabricTransactionAction,error) {
	var err error
	var tx *peer.Transaction
	if tx, err = utils.GetTransaction(payload.Data); err != nil {
		return nil,err
	}
	actions := make([]FabricTransactionAction,len(tx.Actions))

	for i, action := range tx.Actions {
		act := FabricTransactionAction{}

		var ca        *peer.ChaincodeAction
		var capayload *peer.ChaincodeActionPayload
		capayload, ca, err = utils.GetPayloads(action)
		if  err != nil {
			return nil,err
		}
		endorsers := make([]string,len(capayload.Action.Endorsements))
		for j, endorser := range capayload.Action.Endorsements {
			var  subject string
			subject, err = decodeSerializedIdentity(endorser.Endorser)
			if err != nil {
				return nil,err
			}
			endorsers[j] = subject
		}
		act.Endorsers = endorsers

		txRWSet := &rwsetutil.TxRwSet{}
		err = txRWSet.FromProtoBytes(ca.Results)
		if err != nil {
			return nil,err
		}

		rwSets := make([]FabricTransactionActionRWSet,0)
		for _, nsRWSet := range txRWSet.NsRwSets {
			ns := nsRWSet.NameSpace
			if ns != "lscc" {   // skip system chaincode
				rwSet := FabricTransactionActionRWSet{}
				rwSet.Cc = ns
				rset := make([]string,len(nsRWSet.KvRwSet.Reads))
				for i, kvRead := range nsRWSet.KvRwSet.Reads {
					rset[i] = fmt.Sprintf("key=%v,version=%v",kvRead.Key,kvRead.Version)
				}
				rwSet.RSet = rset

				wset := make([]string,len(nsRWSet.KvRwSet.Writes))
				for i, kvWrite := range nsRWSet.KvRwSet.Writes {
					wset[i] = fmt.Sprintf("key=%v,isDelete=%v,value=%v",kvWrite.Key,kvWrite.IsDelete,string(kvWrite.Value))
				}
				rwSet.WSet = wset
				rwSets = append(rwSets,rwSet)
			}
		}
		act.RWSet = rwSets
		actions[i]= act
	}
	return actions,nil
}

func parseConfig(payload *common.Payload) ([]byte,error) {
	var err error

	var configEnvelope *common.ConfigEnvelope
	configEnvelope, err = configtx.UnmarshalConfigEnvelope(payload.Data)
	if  err != nil {
		return nil,err
	}
	config := configEnvelope.Config
	conf ,err := json.Marshal(config)
	if  err != nil {
		return nil,err
	}
	return conf,nil
}