// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package store

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// CreateFeed is the `createFeed` instruction.
type CreateFeed struct {
	Description *string
	Decimals    *uint8
	Granularity *uint8
	LiveLength  *uint32

	// [0] = [] store
	//
	// [1] = [WRITE] feed
	//
	// [2] = [SIGNER] authority
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateFeedInstructionBuilder creates a new `CreateFeed` instruction builder.
func NewCreateFeedInstructionBuilder() *CreateFeed {
	nd := &CreateFeed{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetDescription sets the "description" parameter.
func (inst *CreateFeed) SetDescription(description string) *CreateFeed {
	inst.Description = &description
	return inst
}

// SetDecimals sets the "decimals" parameter.
func (inst *CreateFeed) SetDecimals(decimals uint8) *CreateFeed {
	inst.Decimals = &decimals
	return inst
}

// SetGranularity sets the "granularity" parameter.
func (inst *CreateFeed) SetGranularity(granularity uint8) *CreateFeed {
	inst.Granularity = &granularity
	return inst
}

// SetLiveLength sets the "liveLength" parameter.
func (inst *CreateFeed) SetLiveLength(liveLength uint32) *CreateFeed {
	inst.LiveLength = &liveLength
	return inst
}

// SetStoreAccount sets the "store" account.
func (inst *CreateFeed) SetStoreAccount(store ag_solanago.PublicKey) *CreateFeed {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(store)
	return inst
}

// GetStoreAccount gets the "store" account.
func (inst *CreateFeed) GetStoreAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// SetFeedAccount sets the "feed" account.
func (inst *CreateFeed) SetFeedAccount(feed ag_solanago.PublicKey) *CreateFeed {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(feed).WRITE()
	return inst
}

// GetFeedAccount gets the "feed" account.
func (inst *CreateFeed) GetFeedAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// SetAuthorityAccount sets the "authority" account.
func (inst *CreateFeed) SetAuthorityAccount(authority ag_solanago.PublicKey) *CreateFeed {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *CreateFeed) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst CreateFeed) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_CreateFeed,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst CreateFeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CreateFeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Description == nil {
			return errors.New("Description parameter is not set")
		}
		if inst.Decimals == nil {
			return errors.New("Decimals parameter is not set")
		}
		if inst.Granularity == nil {
			return errors.New("Granularity parameter is not set")
		}
		if inst.LiveLength == nil {
			return errors.New("LiveLength parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Store is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Feed is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Authority is not set")
		}
	}
	return nil
}

func (inst *CreateFeed) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("CreateFeed")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=4]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Description", *inst.Description))
						paramsBranch.Child(ag_format.Param("   Decimals", *inst.Decimals))
						paramsBranch.Child(ag_format.Param("Granularity", *inst.Granularity))
						paramsBranch.Child(ag_format.Param(" LiveLength", *inst.LiveLength))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("    store", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("     feed", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("authority", inst.AccountMetaSlice[2]))
					})
				})
		})
}

func (obj CreateFeed) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Description` param:
	err = encoder.Encode(obj.Description)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(obj.Decimals)
	if err != nil {
		return err
	}
	// Serialize `Granularity` param:
	err = encoder.Encode(obj.Granularity)
	if err != nil {
		return err
	}
	// Serialize `LiveLength` param:
	err = encoder.Encode(obj.LiveLength)
	if err != nil {
		return err
	}
	return nil
}
func (obj *CreateFeed) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Description`:
	err = decoder.Decode(&obj.Description)
	if err != nil {
		return err
	}
	// Deserialize `Decimals`:
	err = decoder.Decode(&obj.Decimals)
	if err != nil {
		return err
	}
	// Deserialize `Granularity`:
	err = decoder.Decode(&obj.Granularity)
	if err != nil {
		return err
	}
	// Deserialize `LiveLength`:
	err = decoder.Decode(&obj.LiveLength)
	if err != nil {
		return err
	}
	return nil
}

// NewCreateFeedInstruction declares a new CreateFeed instruction with the provided parameters and accounts.
func NewCreateFeedInstruction(
	// Parameters:
	description string,
	decimals uint8,
	granularity uint8,
	liveLength uint32,
	// Accounts:
	store ag_solanago.PublicKey,
	feed ag_solanago.PublicKey,
	authority ag_solanago.PublicKey) *CreateFeed {
	return NewCreateFeedInstructionBuilder().
		SetDescription(description).
		SetDecimals(decimals).
		SetGranularity(granularity).
		SetLiveLength(liveLength).
		SetStoreAccount(store).
		SetFeedAccount(feed).
		SetAuthorityAccount(authority)
}
