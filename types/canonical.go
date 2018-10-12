package types

import (
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"
)

// Canonical* wraps the structs in types for amino encoding them for use in SignBytes / the Signable interface.

// TimeFormat is used for generating the sigs
const TimeFormat = time.RFC3339Nano

type CanonicalBlockID struct {
	Hash        cmn.HexBytes
	PartsHeader CanonicalPartSetHeader
}

type CanonicalPartSetHeader struct {
	Hash  cmn.HexBytes
	Total int
}

type CanonicalProposal struct {
	Version          uint64        `binary:"fixed64"`
	Height           int64         `binary:"fixed64"`
	Round            int64         `binary:"fixed64"`
	Type             SignedMsgType // type alias for byte
	POLRound         int64         `binary:"fixed64"`
	Timestamp        time.Time
	BlockPartsHeader CanonicalPartSetHeader
	POLBlockID       CanonicalBlockID
	ChainID          string
}

type CanonicalVote struct {
	Version   uint64        `binary:"fixed64"`
	Height    int64         `binary:"fixed64"`
	Round     int64         `binary:"fixed64"`
	Type      SignedMsgType // type alias for byte
	Timestamp time.Time
	BlockID   CanonicalBlockID
	ChainID   string
}

type CanonicalHeartbeat struct {
	Version          uint64 `binary:"fixed64"`
	Height           int64  `binary:"fixed64"`
	Round            int    `binary:"fixed64"`
	Type             byte
	Sequence         int `binary:"fixed64"`
	ValidatorAddress Address
	ValidatorIndex   int
	ChainID          string
}

//-----------------------------------
// Canonicalize the structs

func CanonicalizeBlockID(blockID BlockID) CanonicalBlockID {
	return CanonicalBlockID{
		Hash:        blockID.Hash,
		PartsHeader: CanonicalizePartSetHeader(blockID.PartsHeader),
	}
}

func CanonicalizePartSetHeader(psh PartSetHeader) CanonicalPartSetHeader {
	return CanonicalPartSetHeader{
		psh.Hash,
		psh.Total,
	}
}

func CanonicalizeProposal(chainID string, proposal *Proposal) CanonicalProposal {
	return CanonicalProposal{
		Version:          0, // TODO
		Height:           proposal.Height,
		Round:            int64(proposal.Round), // cast int->int64 to make amino encode it fixed64 (does not work for int)
		Type:             ProposalType,
		POLRound:         int64(proposal.POLRound),
		Timestamp:        proposal.Timestamp,
		BlockPartsHeader: CanonicalizePartSetHeader(proposal.BlockPartsHeader),
		POLBlockID:       CanonicalizeBlockID(proposal.POLBlockID),
		ChainID:          chainID,
	}
}

func CanonicalizeVote(chainID string, vote *Vote) CanonicalVote {
	return CanonicalVote{
		Version:   0, // TODO
		Height:    vote.Height,
		Round:     int64(vote.Round), // cast int->int64 to make amino encode it fixed64 (does not work for int)
		Type:      vote.Type,
		Timestamp: vote.Timestamp,
		BlockID:   CanonicalizeBlockID(vote.BlockID),
		ChainID:   chainID,
	}
}

func CanonicalizeHeartbeat(chainID string, heartbeat *Heartbeat) CanonicalHeartbeat {
	return CanonicalHeartbeat{
		Version:          0, // TODO
		Height:           heartbeat.Height,
		Round:            heartbeat.Round,
		Type:             byte(HeartbeatType),
		Sequence:         heartbeat.Sequence,
		ValidatorAddress: heartbeat.ValidatorAddress,
		ValidatorIndex:   heartbeat.ValidatorIndex,
		ChainID:          chainID,
	}
}

// CanonicalTime can be used to stringify time in a canonical way.
func CanonicalTime(t time.Time) string {
	// Note that sending time over amino resets it to
	// local time, we need to force UTC here, so the
	// signatures match
	return tmtime.Canonical(t).Format(TimeFormat)
}
